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

var alls = make(map[string]Tasks)

func SchedMyn(task *WorkerRequest, worker *WorkerHandle) bool {
	taskid := task.Sector.ID.Number.String()
	tasktype := task.TaskType.Short()
	workername := worker.Info.Hostname
	log.Debugf("taskid is %s woker", taskid, workername)
	lck.Lock()
	defer lck.Unlock()
	if worker.Info.Hostname != "miner" {
		if _, ok := all[workername]; !ok {
			alls[workername] = Tasks{Tasklist: make(map[string]string)}
		}
		if _, ok := alls[workername].Tasklist[taskid]; ok {
			if all[workername].Tasklist[taskid] == "FIN" {
				delete(all[workername].Tasklist, taskid)
				return true
			}
			alls[workername].Tasklist[taskid] = tasktype
			log.Debugf("update tasktype is %s woker", taskid, workername, tasktype)
			return true
		}

		if tasktype == "AP" && len(all[workername].Tasklist) < 5 {
			all[workername].Tasklist[taskid] = tasktype
			log.Debugf("add taskid for %s woker", taskid, workername)
			return true
		}
	}
	log.Debugf("alls is %s ", alls)
	return false
}
func wAllworkersToJson() error {
	lock.Lock()
	defer lock.Unlock()
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
