package sealer

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
)

//type Local struct {
//	ServerName string   `json:"serverName"`
//	Sectors    []string `json:"sectors"`
//}

//var ch = make(chan string)

// var sectors []string
// var ws = []Local{}
var mutex sync.Mutex

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {
	//stat =
	//per(task.TaskType.Short(), task.Sector.ID.Number.String(), worker.Info.Hostname)
	//log.Debugf(task.Sector.ID.Number.String())
	//per(sector task.Sector.ID.Number.String())

	//wrpc, _ := worker.workerRpc.TaskTypes(task.Ctx)
	//for taskType, s := range wrpc {
	//	fmt.Println("rpc"+taskType.Short(), s)
	//}
	//for id, handle := range workers {
	//	id.String()
	//	log.Debugf(handle.Info.Hostname)
	//}
	//for taskType, i := range worker.active.taskCounters {
	//log.Debugf(worker.Info.Hostname)
	/*if len(ws) == 0 {
		neww := &Local{
			ServerName: worker.Info.Hostname,
			Sectors:    []string{"1"},
		}
		ws = append(ws, *neww)
		log.Debugf("新建了" + neww.ServerName)
	}
	for _, w := range ws {
		if w.ServerName != worker.Info.Hostname {
			neww := &Local{
				ServerName: worker.Info.Hostname,
				Sectors:    []string{"1"},
			}
			ws = append(ws, *neww)
			log.Debugf("新建了" + w.ServerName)
		}
	}
	//log.Debugf(len(ws))
	//log.Debugf("taskCounters"+taskType.Short(), i)
	//log.Debugf("task is ----------------------" + task.TaskType.Short())
	for _, w := range ws {
		ss := w.Sectors
		for _, sector := range ss {
			if sector == task.Sector.ID.Number.String() {
				if task.TaskType.Short() == "FIN" {
					ss = ss[0 : len(w.Sectors)-1]
					log.Debugf("拉取文件中")

					return true
				}
				log.Debugf(w.ServerName + "正在执行" + task.TaskType.Short())
				return true
			} else {
				log.Debugf(sector + "is" + task.Sector.ID.Number.String())
				//if task.TaskType.Short() == "AP" {
				w.Sectors = append(w.Sectors, task.Sector.ID.Number.String())

				log.Debugf("分配了" + task.Sector.ID.Number.String() + "给" + worker.Info.Hostname)
				return true
				//}
			}
		}
	}
	/*if task.TaskType.Short() == "FIN" && len(se) > 0 {
		delete(se, task.TaskType.Short())
		log.Debugf("拉取文件中")
		return true
		//se[task.Sector.ID.Number.String()] = worker.Info.Hostname
	}
	*/

	h, ok := scene.Load(task.Sector.ID.Number.String())
	if ok {
		if task.TaskType.Short() == "FIN" && h == worker.Info.Hostname {
			scene.Delete(task.Sector.ID.Number.String())
			log.Debugf("拉取文件中")
			return true
		}
		if h == worker.Info.Hostname {

			log.Debugf(worker.Info.Hostname + "正在执行" + task.Sector.ID.Number.String() + "----" + task.TaskType.Short())

			return true

		}

	}
	if !ok && task.TaskType.Short() == "AP" && worker.Info.Hostname != "hcxj-10-0-1-185" {
		scene.Store(task.Sector.ID.Number.String(), worker.Info.Hostname)
		log.Debugf("分配了AP" + task.Sector.ID.Number.String() + "给" + worker.Info.Hostname)
		write()
		return true
	}

	//log.Debugf("task is ----------------------" + task.TaskType.Short())
	//if worker.Info.Hostname == "hcxj-10-0-1-185" {
	//	return true
	//}
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

func read() {
	//os.ReadFile("/home/ts/json")
	f, err := os.ReadFile("/home/ts/json")
	if err != nil {
		panic(err)
	}

	//var tmpMap map[string]interface{}
	if err := json.Unmarshal(f, &scene); err != nil {
		panic(err)
	}
	//for key, value := range tmpMap {
	//	scene.Store(key, value)
	//}
	log.Debugf("读取完成")

}

func write() {
	mutex.Lock()
	/*

		log.Debugf("LOCK")
		whitelist := map[string]string{}
		scene.Range(func(k, v interface{}) bool {
			whitelist[fmt.Sprint(k)] = fmt.Sprint(v)
			return true
		})
		fmt.Println("Done.")

	*/
	f, err := os.OpenFile("/home/ts/json", os.O_WRONLY|os.O_CREATE, 0666)

	//syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(scene)
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(buf.Bytes()))
	if err != nil {
		panic(err)
	}
	//syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	mutex.Unlock()
	log.Debugf("UNLOCK")
}
