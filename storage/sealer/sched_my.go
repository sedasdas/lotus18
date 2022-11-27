package sealer

func SchedLocal(task *WorkerRequest, request *SchedWindowRequest, worker *WorkerHandle) bool {
	log.Debugf(worker.Info.Hostname)
	log.Debugf(task.Sector.ID.Number.String())
	return false
}
