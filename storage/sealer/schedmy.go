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

func SchedMy(task *WorkerRequest, worker *WorkerHandle) bool {
	taskid := task.Sector.ID.Number.String()
	log.Debugf("taskid is %s woker", taskid, worker.Info.Hostname)
	if worker.Info.Hostname != "miner" {
		findWorkertoAllworkers(worker.Info.Hostname)
		for _, w := range allworkers {
			if w.getWorker() == worker.Info.Hostname {
				if status, _ := w.getTask(taskid); status == "FIN" {
					if err := w.delTask(taskid); err != nil {
						return false
					}
					return true
				}
				if w.getTaskListLen() < 4 && worker.Info.Hostname != "miner" {
					if err := w.addTask(taskid, task.TaskType.Short(), worker.Info.Hostname); err != nil {

						return false
					}
					log.Debugf("add task %s to worker %s do   %s", taskid, worker.Info.Hostname, task.TaskType.Short())
					return true
				}
				return true
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
		tasklist: sync.Map{},
		name:     name,
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
	file, _ := os.Create(".workers.json")
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
