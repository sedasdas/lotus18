package sealer

type local struct {
	sectors []string
}

var sectors []string

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {
	log.Debugf(worker.Info.Hostname)
	log.Debugf(task.Sector.ID.Number.String())
	//per(sector task.Sector.ID.Number.String())

	if worker.Info.Hostname == "hcxj-10-0-5-71" {
		log.Debugf("RIGHT??????????????????????????????")
		sectors = append(sectors, task.Sector.ID.Number.String())
		log.Debugf("len=%d cap=%d slice=%v\n", len(sectors), cap(sectors), sectors)
		return true
	}

	return false
}
func read(l local) {

}
func per(sector string, l local) {

}
