package main

import (
	"fmt"
	"os"

	"github.com/gdata/customer-saga/workers"
)

func init() {
	dir, _ := os.Getwd()
	os.Setenv("APP_WATCH_PATH", fmt.Sprintf("%v/%v", dir, ".files/")) // "/mnt/c/Projects/.samples/GData-Pipeline/customer-saga/.files/"
}

func main() {
	Create().
		RegisterWorker(workers.CreateWorkerFactory[*workers.CsvProcessingWorker]()).
		Run()
}
