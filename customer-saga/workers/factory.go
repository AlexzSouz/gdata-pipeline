package workers

import (
	"github.com/gdata/customer-saga/abstractions"
)

type IWorker interface {
	ExecuteWithContext(ctx abstractions.IAppContext)
}

func CreateWorkerFactory[T IWorker]() T {
	return *new(T)
}
