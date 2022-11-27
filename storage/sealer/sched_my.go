package sealer

type local struct {
	max int
}

func SchedLocal(l local, task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {
	log.Debugf(worker.Info.Hostname)
	log.Debugf(task.Sector.ID.Number.String())
	if worker.Info.Hostname == "hcxj-10-0-5-71" {
		log.Debugf("RIGHT??????????????????????????????")
		return true
	}

	return false
}
