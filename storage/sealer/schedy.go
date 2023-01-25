package sealer

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
)

var lck sync.Mutex

type Tasks struct {
	Tasklist map[string]string `json:"tasklist"`
}

func (t *Tasks) getTaskCountPc1(status string) int {
	p1c := 0
	for _, s := range t.Tasklist {
		if s == status {
			p1c++
		}
	}

	return p1c
}

var alls = make(map[string]*Tasks)

func SchedMyn(task *WorkerRequest, worker *WorkerHandle) bool {
	taskid := task.Sector.ID.Number.String()
	tasktype := task.TaskType.Short()
	workername := worker.Info.Hostname
	//log.Debugf(workername, taskid, tasktype)
	lck.Lock()
	defer lck.Unlock()
	if worker.Info.Hostname != "miner" {
		if _, ok := alls[workername]; !ok {
			alls[workername] = &Tasks{Tasklist: make(map[string]string)}
			log.Debugf("add new worker %s", alls[workername])
		}
		if tasktype == "AP" && alls[workername].getTaskCountPc1("P1C") < 4 {
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("add taskid ap for %s woker", taskid, workername)
			return true
		}

		if _, ok := alls[workername].Tasklist[taskid]; ok {
			if tasktype == "AP" && alls[workername].getTaskCountPc1("P1C") < 4 {
				alls[workername].Tasklist[taskid] = tasktype
				log.Debugf("update taskid ap for %s woker", taskid, workername)
				return true
			}
			log.Debugf("task %s  is added ", taskid)
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("update taskid %s type is  %s  for woker %s ", taskid, tasktype, workername)
			return true
		}

	}
	if tasktype == "FIN" {
		return true
	}

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
