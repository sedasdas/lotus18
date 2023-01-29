package sealer

import (
	"bytes"
	"encoding/json"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
	"os"
	"strings"
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

func SchedMyn(task *WorkerRequest, worker *WorkerHandle, workers map[storiface.WorkerID]*WorkerHandle) bool {
	lck.Lock()
	defer lck.Unlock()
	taskid := task.Sector.ID.Number.String()
	tasktype := task.TaskType.Short()
	workername := worker.Info.Hostname
	basename := workername[:len(workername)-2]
	//log.Debugf(workername, taskid, tasktype)

	if worker.Info.Hostname != "miner" {
		if _, ok := alls[workername]; !ok {
			alls[workername] = &Tasks{Tasklist: make(map[string]string)}
			log.Debugf("add new worker %s", workername)
			log.Debugf("add new worker %s", alls[workername])
		}

		if strings.Contains(workername, "p1") && tasktype == "AP" && alls[workername].getTaskCountPc1("P1C") < 4 && alls[workername].getTaskCountPc1("AP") < 4 {
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("add task ap for %s woker %s", taskid, workername)
			return true
		}

		if _, ok := alls[basename+"p1"].Tasklist[taskid]; ok {
			if tasktype == "AP" || tasktype == "PC1" {
				if strings.Contains(workername, "p1") && alls[workername].getTaskCountPc1("P1C") < 4 {
					alls[workername].Tasklist[taskid] = tasktype
					log.Debugf("update task ap or p1 for %s woker %s", taskid, workername)
					return true
				}
			} else if strings.Contains(workername, "p2") {
				//log.Debugf("workername %s", handle.Info.Hostname)
				alls[workername].Tasklist[taskid] = tasktype
				log.Debugf("update task %s for %s woker %s", tasktype, taskid, workername)
				return true
			}
		}
		//log.Debugf(" %s woker is busy", workernam e)
		return false
	}

	return false
}
func wAllworkersToJson() error {

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
