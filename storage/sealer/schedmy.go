package sealer

import (
	"sync"
)

var lock sync.Mutex

type MyWorker struct {
	tasklist (map[string]string) `json:"tasklist"`
	name     string              `json:"name"`
	sync.RWMutex
}

var allworkers = []MyWorker{}

func (myw *MyWorker) addTask(taskid string, status string, worker string) {
	myw.Lock()
	myw.tasklist[taskid] = status
	myw.Unlock()
	myw.name = worker
}
func (myw *MyWorker) delTask(taskid string) {
	myw.Lock()
	delete(myw.tasklist, taskid)
	myw.Unlock()
}
func (myw *MyWorker) getTask(taskid string) string {
	myw.RLock()
	value := myw.tasklist[taskid]
	myw.RUnlock()
	return value
}
func (myw *MyWorker) getWorker() string {
	return myw.name
}
func (myw *MyWorker) getTaskList() map[string]string {
	return myw.tasklist
}
func (myw *MyWorker) getTaskListLen() int {
	return len(myw.tasklist)
}

func SchedMy(task *WorkerRequest, worker *WorkerHandle) bool {
	for _, w := range allworkers {
		if w.getWorker() == worker.Info.Hostname {
			if w.getTask(task.Sector.ID.Number.String()) == "FIN" {
				w.delTask(task.Sector.ID.Number.String())
				return true
			}
			if w.getTaskListLen() < 4 && worker.Info.Hostname != "miner" {
				w.addTask(task.Sector.ID.Number.String(), task.TaskType.Short(), worker.Info.Hostname)

				return true
			}
			return true
		}
	}
	addWorkertoAllworkers(worker.Info.Hostname, task)
	return true
}
func addWorkertoAllworkers(name string, task *WorkerRequest) {
	lock.Lock()
	defer lock.Unlock()
	wok := &MyWorker{
		tasklist: make(map[string]string),
	}
	wok.name = name
	wok.tasklist[task.Sector.ID.Number.String()] = task.TaskType.Short()
	allworkers = append(allworkers, *wok)
}

func (myw MyWorker) delWorkertoAllworkers() {
	for i, w := range allworkers {
		if w.getWorker() == myw.getWorker() {
			allworkers = append(allworkers[:i], allworkers[i+1:]...)
		}
	}
}
