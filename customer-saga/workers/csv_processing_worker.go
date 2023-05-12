package workers

import (
	"fmt"
	"os"

	"github.com/gdata/customer-saga/abstractions"
	"github.com/gdata/customer-saga/watchers"
	"go.opentelemetry.io/otel/attribute"
)

type CsvProcessingWorker struct {
}

func (w *CsvProcessingWorker) ExecuteWithContext(ctx abstractions.IAppContext) {
	workerType := fmt.Sprintf("%T", w)
	ctx, span := ctx.CreateSpan(workerType)
	defer span.End()

	defer func() {
		if err := recover(); err != nil {
			ctx.Terminate()
		}
	}()

	span.SetAttributes(attribute.String("worker.type", workerType))

	fsWatcher := watchers.CreateFileSystemWatcher(ctx)
	fsWatcher.Watch(os.Getenv("APP_WATCH_PATH")) // TODO : Pass a Command Pattern object to the watcher to process the data
}
