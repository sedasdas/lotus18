package sealer

import (
	"bytes"
	"encoding/json"
	"fmt"
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
var workerTaskCount = make(map[string]int)
var someMapMutex = sync.RWMutex{}

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {

	h, ok := scene.Load(task.Sector.ID.Number.String())
	if ok {
		if task.TaskType.Short() == "FIN" && h == worker.Info.Hostname {
			scene.Delete(task.Sector.ID.Number.String())
			log.Debugf("拉取文件中")
			return true
		}
		if h == worker.Info.Hostname {
			assignTask(worker.Info.Hostname)

			log.Debugf(worker.Info.Hostname + "正在执行" + task.Sector.ID.Number.String() + "----" + task.TaskType.Short())

			return true

		}

	}
	if !ok && task.TaskType.Short() == "AP" && worker.Info.Hostname != "miner" {

		assignTask(worker.Info.Hostname)
		scene.Store(task.Sector.ID.Number.String(), worker.Info.Hostname)
		log.Debugf("分配了AP" + task.Sector.ID.Number.String() + "给" + worker.Info.Hostname)
		write()
		return true
	}

	return false
}
func assignTask(worker string) bool {
	if workerTaskCount[worker] >= 4 {
		log.Debugf("Worker %s already has 4 tasks, cannot assign more", worker)
		return false
	}
	// assign task to worker
	someMapMutex.Lock()
	workerTaskCount[worker]++
	someMapMutex.Unlock()
	log.Debugf("Worker %s tasks is  ", worker, workerTaskCount[worker])
	return true
}
func read() {
	//os.ReadFile("/home/ts/json")
	f, err := os.ReadFile("/home/cali/lotusjson")
	if err != nil {
		panic(err)
	}

	var tmpMap map[string]interface{}
	if err := json.Unmarshal(f, &tmpMap); err != nil {
		panic(err)
	}
	for key, value := range tmpMap {
		scene.Store(key, value)
	}
	log.Debugf("读取完成")

}

func write() {
	mutex.Lock()
	log.Debugf("LOCK")
	whitelist := map[string]string{}
	scene.Range(func(k, v interface{}) bool {
		whitelist[fmt.Sprint(k)] = fmt.Sprint(v)
		return true
	})
	fmt.Println("Done.")
	f, err := os.OpenFile("/home/cali/lotusjson", os.O_WRONLY|os.O_TRUNC, 0666)

	//syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(whitelist)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(buf.Bytes()))
	if err != nil {
		panic(err)
	}
	//syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	mutex.Unlock()
	log.Debugf("写入文件")
}
