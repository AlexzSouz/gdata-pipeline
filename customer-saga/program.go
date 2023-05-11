package main

import (
	"os"

	"github.com/gdata/customer-saga/workers"
)

func main() {
	// Temp initializations
	os.Setenv("APP_WATCH_PATH", "/mnt/c/Projects/.samples/GData-Pipeline/customer-saga/.files/")

	// Startup
	Create().
		RegisterWorker(workers.CreateWorkerFactory[*workers.CsvProcessingWorker]()).
		Run()
}
