package sealer

import (
	"context"
	"github.com/filecoin-project/lotus/storage/sealer/sealtasks"
	"sync"
)

type WindowSelector func(sh *Scheduler, queueLen int, acceptableWindows [][]int, windows []SchedWindow) int

type AssignerCommon struct {
	taskChannel chan *Scheduler
	wg          sync.WaitGroup
}

func NewAssignerCommon() *AssignerCommon {
	a := &AssignerCommon{
		taskChannel: make(chan *Scheduler),
	}
	return a
}

func (a *AssignerCommon) TrySched(sh *Scheduler) {
	windowsLen := len(sh.OpenWindows)
	queueLen := sh.SchedQueue.Len()
	log.Debugf("SCHED %d queued; %d open windows", queueLen, windowsLen)
	if windowsLen == 0 || queueLen == 0 {
		// nothing to schedule on
		return
	}

	acceptableWindowsPool := &sync.Pool{New: func() interface{} { return make([]int, 0, windowsLen) }}

	a.wg.Add(queueLen)
	for i := 0; i < queueLen; i++ {
		go func(i int) {
			defer a.wg.Done()
			task := (*sh.SchedQueue)[i]
			task.IndexHeap = i
			acceptableWindows := acceptableWindowsPool.Get().([]int)[:0]
			var havePreferred bool
			for wnd, windowRequest := range sh.OpenWindows {
				worker, ok := sh.Workers[windowRequest.Worker]
				if !ok {
					log.Errorf("worker referenced by windowRequest not found (worker: %s)", windowRequest.Worker)
					continue
				}

				if !worker.Enabled {
					log.Debugw("skipping disabled worker", "worker", windowRequest.Worker)
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
					acceptableWindows = acceptableWindows[:0]
					havePreferred = true
				}

				acceptableWindows = append(acceptableWindows, wnd)
			}

			if len(acceptableWindows) == 0 {
				return
			}
			acceptableWindowsPool.Put(acceptableWindows)
		}(i)
	}
	a.wg.Wait()

}
