package sealer

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

type MyWorker struct {
	tasklist sync.Map
	name     string `json:"name"`
}

var allworkers = []MyWorker{}

func (myw *MyWorker) addTask(taskid string, status string, worker string) error {
	if _, loaded := myw.tasklist.LoadOrStore(taskid, status); loaded {
		return fmt.Errorf("task with id %s already exists", taskid)
	}
	myw.name = worker
	return nil
}
func (myw *MyWorker) delTask(taskid string) error {
	if _, ok := myw.tasklist.Load(taskid); !ok {
		return fmt.Errorf("task with id %s not found", taskid)
	}
	myw.tasklist.Delete(taskid)
	return nil
}
func (myw *MyWorker) getTask(taskid string) (string, error) {
	if value, ok := myw.tasklist.Load(taskid); !ok {
		return "", fmt.Errorf("task with id %s not found", taskid)
	} else {
		return value.(string), nil
	}
}
func (myw *MyWorker) getWorker() string {
	return myw.name
}
func (myw *MyWorker) getTaskListLen() int {
	var length int
	myw.tasklist.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func SchedMy(task *WorkerRequest, worker *WorkerHandle) (bool, error) {
	taskid := task.Sector.ID.Number.String()
	for _, w := range allworkers {
		if w.getWorker() == worker.Info.Hostname {
			if status, _ := w.getTask(taskid); status == "FIN" {
				if err := w.delTask(taskid); err != nil {
					return false, err
				}
				return true, nil
			}
			if w.getTaskListLen() < 4 && worker.Info.Hostname != "miner" {
				if err := w.addTask(taskid, task.TaskType.Short(), worker.Info.Hostname); err != nil {
					return false, err
				}
				return true, nil
			}
			return true, nil
		}
	}
	if err := addWorkertoAllworkers(worker.Info.Hostname, task); err != nil {
		return false, err
	}
	return true, nil
}
func addWorkertoAllworkers(name string, task *WorkerRequest) error {
	lock.Lock()
	defer lock.Unlock()
	wok := &MyWorker{
		tasklist: sync.Map{},
		name:     name,
	}
	wok.tasklist.Store(task.Sector.ID.Number.String(), task.TaskType.Short())
	log.Debugf("分配了" + task.Sector.ID.Number.String() + task.TaskType.Short() + "给" + name)
	allworkers = append(allworkers, *wok)
	log.Debugf("add worker %s to allworkers", name)
	return nil
}
