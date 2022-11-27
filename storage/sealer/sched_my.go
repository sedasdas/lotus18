package sealer

type Local struct {
	ServerName string    `json:"serverName"`
	Sectors    [0]string `json:"sectors"`
}

var s = 3

// var sectors []string
var sectors = []string{"3", "4", "5"}

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {
	log.Debugf(worker.Info.Hostname)
	log.Debugf(task.Sector.ID.Number.String())
	//per(sector task.Sector.ID.Number.String())

	if worker.Info.Hostname == "hcxj-10-0-5-71" {
		log.Debugf("RIGHT??????????????????????????????")
		//sectors = append(sectors, task.Sector.ID.Number.String())
		for _, sector := range sectors {
			log.Debugf(sector)
			if task.Sector.ID.Number.String() == sector {
				log.Debugf("分配了" + sector)
				return true
			}
		}
		//log.Debugf("len=%d cap=%d slice=%v\n", len(sectors), cap(sectors), sectors)
	}

	return false
}
