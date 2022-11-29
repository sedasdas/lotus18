package sealer

import (
	"fmt"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
)

type Local struct {
	ServerName string   `json:"serverName"`
	Sectors    []string `json:"sectors"`
}

var stat storiface.WorkerStats
var se = make(map[string]string)

//var ch = make(chan string)

// var sectors []string
var sectors []string

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle, workers map[storiface.WorkerID]*WorkerHandle) bool {
	//stat =
	//per(task.TaskType.Short(), task.Sector.ID.Number.String(), worker.Info.Hostname)
	//log.Debugf(task.Sector.ID.Number.String())
	//per(sector task.Sector.ID.Number.String())

	wrpc, _ := worker.workerRpc.TaskTypes(task.Ctx)
	for taskType, s := range wrpc {
		fmt.Println("rpc"+taskType.Short(), s)
	}
	for id, handle := range workers {
		id.String()
		log.Debugf(handle.Info.Hostname)
	}
	//for taskType, i := range worker.active.taskCounters {
	log.Debugf(string(len(se)))
	//log.Debugf("taskCounters"+taskType.Short(), i)
	log.Debugf("task is ----------------------" + task.TaskType.Short())
	if task.TaskType.Short() == "FIN" && len(se) > 0 {
		delete(se, task.TaskType.Short())
		return true
		//se[task.Sector.ID.Number.String()] = worker.Info.Hostname
	}

	_, ok := se[task.Sector.ID.Number.String()]
	if ok {
		return true
	}
	if task.TaskType.Short() == "AP" {
		se[task.Sector.ID.Number.String()] = worker.Info.Hostname
		return true
	}

	log.Debugf("task is ----------------------" + task.TaskType.Short())
	//if worker.Info.Hostname == "hcxj-10-0-1-185" {
	//ch <- task.Sector.ID.Number.String()
	//log.Debugf("RIGHT??????????????????????????????")
	//sectors = append(sectors, task.Sector.ID.Number.String())
	//for _, sector := range sectors {
	//log.Debugf(sector)
	//	return true
	//}
	//log.Debugf("len=%d cap=%d slice=%v\n", len(sectors), cap(sectors), sectors)
	//}

	return false
}
