package watchers

import (
	"errors"
	"io/fs"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gdata/customer-saga/abstractions"
	"github.com/gdata/customer-saga/domain"
)

type IFileSystemDirectoryWatcher interface {
	Watch(path string)
	isInvalidDirectory(path string) bool
}

type FileSystemDirectoryWatcher struct {
	DirectoryWatcher
}

func CreateFileSystemWatcher(ctx abstractions.IAppContext) IFileSystemDirectoryWatcher {
	return &FileSystemDirectoryWatcher{
		DirectoryWatcher: DirectoryWatcher{
			Context: ctx,
		},
	}
}

func (w *FileSystemDirectoryWatcher) Watch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		w.Context.Logger().Print("Failed to create a Watcher", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				w.Context.Logger().Printf("Event received for '%s' through operation '%s'", event.Name, event.Op)

				switch event.Op {
				case fsnotify.Create:
					w.Context.Logger().Printf("Starting to process file %q", event.Name)

					service := domain.CustomerService{}
					service.Process(w.Context, event.Name)

					break
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				w.Context.Logger().Fatal("Failed to set watcher.", err)
			}
		}
	}()

	if w.isInvalidDirectory(path) {
		close(done)
		return
	}

	err = watcher.Add(path)
	if err != nil {
		w.Context.Logger().Fatalf("Failed to configure watcher for path at '%v'. %q", path, err)
	}

	<-done
}

func (w *FileSystemDirectoryWatcher) isInvalidDirectory(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		w.Context.Logger().Fatal("Invalid directory provided.", err)
		return errors.Is(err, fs.ErrNotExist)
	}

	return !info.IsDir()
}
