package sealer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var lock sync.Mutex

type MyWorker struct {
	Tasklist map[string]string `json:"tasklist"`
	Name     string            `json:"name"`
}

var allworkers = []MyWorker{}

func (myw *MyWorker) addTask(taskid string, status string, worker string) error {
	myw.Tasklist[taskid] = status
	return nil
}
func (myw *MyWorker) delTask(taskid string) error {
	if _, ok := myw.Tasklist[taskid]; !ok {
		return fmt.Errorf("task with id %s not found", taskid)
	}
	delete(myw.Tasklist, taskid)
	return nil
}
func (myw *MyWorker) getTask(taskid string) string {
	if _, ok := myw.Tasklist[taskid]; ok {
		log.Debugf("task %s  founed", taskid)
		return taskid
	}
	return ""

}
func (myw *MyWorker) getTaskStatus(taskid string) string {
	if status, ok := myw.Tasklist[taskid]; ok {
		return status
	}
	return ""
}
func (myw *MyWorker) updateTaskStatus(taskid string, status string) {
	myw.Tasklist[taskid] = status

}
func (myw *MyWorker) getWorker() string {
	return myw.Name
}
func (myw *MyWorker) getTaskListLen() int {
	return len(myw.Tasklist)
}

func SchedMy(task *WorkerRequest, worker *WorkerHandle) bool {
	taskid := task.Sector.ID.Number.String()
	log.Debugf("taskid is %s woker", taskid, worker.Info.Hostname)
	if worker.Info.Hostname != "miner" {
		findWorkertoAllworkers(worker.Info.Hostname)
		for _, w := range allworkers {
			if w.getWorker() == worker.Info.Hostname {
				if status := w.getTaskStatus(taskid); status == "FIN" {
					if err := w.delTask(taskid); err != nil {
						return false
					}
					return true
				}
				if w.getTaskListLen() < 4 {
					w.addTask(taskid, task.TaskType.Short(), worker.Info.Hostname)
					log.Debugf("add task %s to worker %s do   %s", taskid, worker.Info.Hostname, task.TaskType.Short())
					return true
				}
				if w.getTask(taskid) == taskid {
					w.updateTaskStatus(taskid, task.TaskType.Short())
					log.Debugf("update task %s to worker %s do   %s", taskid, worker.Info.Hostname, task.TaskType.Short())
				}
				log.Debugf("worker %s is busy tasklen is  %s  ", worker.Info.Hostname, w.getTaskListLen())
				return false
			}

		}
	}
	writeAllworkersToJson()
	return false
}
func addWorkertoAllworkers(name string) error {
	lock.Lock()
	defer lock.Unlock()
	wok := &MyWorker{
		Tasklist: make(map[string]string),
		Name:     name,
	}
	allworkers = append(allworkers, *wok)
	log.Debugf("add worker %s to allworkers", name)
	return nil
}
func findWorkertoAllworkers(wname string) {

	for _, w := range allworkers {
		if w.getWorker() == wname {
			log.Debugf("find worker %s in allworkers", wname)
			return
		}
	}
	addWorkertoAllworkers(wname)
}
func writeAllworkersToJson() error {
	lock.Lock()
	defer lock.Unlock()
	file, _ := os.OpenFile(".workers.json", os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	log.Debugf("write allworkers to .workers.json", allworkers)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(allworkers)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(buf.Bytes()))
	if err != nil {
		panic(err)
	}
	return nil
}
func readAllworkersFromJson() ([]MyWorker, error) {
	file, err := os.ReadFile(".workers.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(file, &allworkers); err != nil {
		panic(err)
	}
	log.Debugf("read allworkers from .workers.json", allworkers)
	return allworkers, nil
}
