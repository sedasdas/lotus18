package sealer

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
)

var lck sync.Mutex

type Tasks struct {
	Tasklist map[string]string `json:"tasklist"`
	P1Count  int
}

var p1c int32

func (t *Tasks) getTaskCountPc1(status string) int {
	p1c := atomic.AddInt32(&p1c, 1)
	for _, s := range t.Tasklist {
		if s == status {
			p1c++
		}
	}
	log.Debugf("getStatus  %s  TaskCount is %s ", status, p1c)
	return int(p1c)

}

var alls = make(map[string]*Tasks)

func SchedMyn(task *WorkerRequest, worker *WorkerHandle) bool {
	taskid := task.Sector.ID.Number.String()
	tasktype := task.TaskType.Short()
	workername := worker.Info.Hostname
	log.Debugf("taskid is %s woker", taskid, workername)
	lck.Lock()
	defer lck.Unlock()
	if worker.Info.Hostname != "miner" {
		if _, ok := alls[workername]; !ok {
			alls[workername] = &Tasks{Tasklist: make(map[string]string)}
			log.Debugf("add new worker %s", alls[workername])
			log.Debugf("worker tasklist  is %s", alls[workername])
		}

		if tasktype == "AP" && len(alls[workername].Tasklist) < 10 && alls[workername].getTaskCountPc1("PC1") < 3 && alls[workername].getTaskCountPc1("AP") < 3 {
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("add taskid ap for %s woker", taskid, workername)
			return true
		}

		if _, ok := alls[workername].Tasklist[taskid]; ok {
			if tasktype == "FIN" {
				delete(alls[workername].Tasklist, taskid)
				log.Debugf("delete taskid for %s woker", taskid, workername)
				return true
			}
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("update tasktype is %s woker", taskid, workername, tasktype)
			return true
		}
	}
	log.Debugf("alls is %s ", alls)
	//wAllworkersToJson()
	return false
}
func wAllworkersToJson() error {
	lck.Lock()
	defer lck.Unlock()
	file, _ := os.OpenFile(".workers.json", os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	log.Debugf("write allworkers to .workers.json", alls)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(alls)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(buf.Bytes()))
	if err != nil {
		panic(err)
	}
	return nil
}
func rAllworkersFromJson() error {
	file, err := os.ReadFile(".workers.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(file, &alls); err != nil {
		panic(err)
	}
	log.Debugf("read allworkers from .workers.json", alls)
	return nil
}
