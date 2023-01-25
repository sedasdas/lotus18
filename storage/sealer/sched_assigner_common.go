package sealer

import (
	"context"
	"github.com/filecoin-project/lotus/storage/sealer/sealtasks"
	"sync"
)

type WindowSelector func(sh *Scheduler, queueLen int, acceptableWindows [][]int, windows []SchedWindow) int

// AssignerCommon is a task assigner with customizable parts
type AssignerCommon struct {
	WindowSel WindowSelector
}

var _ Assigner = &AssignerCommon{}

func (a *AssignerCommon) TrySched(sh *Scheduler) {
	windowsLen := len(sh.OpenWindows)
	queueLen := sh.SchedQueue.Len()

	log.Debugf("SCHED %d queued; %d open windows", queueLen, windowsLen)

	if windowsLen == 0 || queueLen == 0 {
		// nothing to schedule on
		return
	}

	windows := make([]SchedWindow, windowsLen)
	for i := range windows {
		windows[i].Allocated = *NewActiveResources()
	}
	acceptableWindows := make([][]int, queueLen) // QueueIndex -> []OpenWindowIndex

	// Step 1
	throttle := make(chan struct{}, windowsLen)

	var wg sync.WaitGroup
	wg.Add(queueLen)
	for i := 0; i < queueLen; i++ {
		throttle <- struct{}{}

		go func(sqi int) {
			defer wg.Done()
			defer func() {
				<-throttle
			}()

			task := (*sh.SchedQueue)[sqi]
			task.IndexHeap = sqi

			var havePreferred bool

			for wnd, windowRequest := range sh.OpenWindows {

				worker, ok := sh.Workers[windowRequest.Worker]
				if !ok {
					log.Errorf("worker referenced by windowRequest not found (worker: %s)", windowRequest.Worker)
					// TODO: How to move forward here?
					continue
				}

				if !worker.Enabled {
					log.Debugw("skipping disabled worker", "worker", windowRequest.Worker)
					continue
				}

				needRes := worker.Info.Resources.ResourceSpec(task.Sector.ProofType, task.TaskType)

				// TODO: allow bigger windows
				if !windows[wnd].Allocated.CanHandleRequest(task.SealTask(), needRes, windowRequest.Worker, "schedAcceptable", worker.Info) {
					continue
				}
				if task.TaskType != sealtasks.TTFetch && !SchedMyn(task, worker) {
					continue
				}

				rpcCtx, cancel := context.WithTimeout(task.Ctx, SelectorTimeout)

				ok, preferred, err := task.Sel.Ok(rpcCtx, task.TaskType, task.Sector.ProofType, worker)
				cancel()
				if err != nil {
					log.Errorf("trySched(1) req.Sel.Ok error: %+v", err)
					continue
				}

				if !ok {
					continue
				}

				if havePreferred && !preferred {
					// we have a way better worker for this task
					continue
				}

				if preferred && !havePreferred {
					// all workers we considered previously are much worse choice
					acceptableWindows[sqi] = acceptableWindows[sqi][:0]
					havePreferred = true
				}

				acceptableWindows[sqi] = append(acceptableWindows[sqi], wnd)
			}

			if len(acceptableWindows[sqi]) == 0 {
				return
			}

			// Pick best worker (shuffle in case some workers are equally as good)

		}(i)
	}

	wg.Wait()
	log.Debugf("SCHED windows: %+v", windows)
	log.Debugf("SCHED Acceptable win: %+v", acceptableWindows)

	// Step 2
	scheduled := a.WindowSel(sh, queueLen, acceptableWindows, windows)

	// Step 3

	if scheduled == 0 {
		return
	}

	scheduledWindows := map[int]struct{}{}
	for wnd, window := range windows {
		if len(window.Todo) == 0 {
			// Nothing scheduled here, keep the window open
			continue
		}

		scheduledWindows[wnd] = struct{}{}

		window := window // copy
		select {
		case sh.OpenWindows[wnd].Done <- &window:
		default:
			log.Error("expected sh.OpenWindows[wnd].Done to be buffered")
		}
	}

	// Rewrite sh.OpenWindows array, removing scheduled windows
	newOpenWindows := make([]*SchedWindowRequest, 0, windowsLen-len(scheduledWindows))
	for wnd, window := range sh.OpenWindows {
		if _, scheduled := scheduledWindows[wnd]; scheduled {
			// keep unscheduled windows open
			continue
		}

		newOpenWindows = append(newOpenWindows, window)
	}

	sh.OpenWindows = newOpenWindows
	wAllworkersToJson()
}
