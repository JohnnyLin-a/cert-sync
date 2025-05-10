package functions

import (
	"encoding/json"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/johnnylin-a/cert-sync/internal/apis"
)

func WatchFiles(filepaths []string, fn func(filename string)) chan struct{} {

	fps, err := json.Marshal(filepaths)
	if err != nil {
		log.Fatal(err)
	}
	apis.LogAndSendNotification("watching files " + string(fps))

	// Create new file watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer func() {
			apis.LogAndSendNotification("watcher stopped")
		}()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					fn(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, filepath := range filepaths {
		err = watcher.Add(filepath)
		if err != nil {
			apis.LogAndSendNotification("failed to watch file " + filepath)
			log.Fatal(err)
		}
	}

	gracefullyStop := make(chan struct{}, 2)
	go func() {
		<-gracefullyStop
		watcher.Close()
		gracefullyStop <- struct{}{}
	}()
	return gracefullyStop
}
