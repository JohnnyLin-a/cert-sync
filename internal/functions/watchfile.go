package functions

import (
	"encoding/json"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/johnnylin-a/cert-sync/internal/apis"
)

func WatchFiles(filepaths []string, fn func(filename string)) chan struct{} {
	notificationSender := apis.GetNotificationsSender()
	fps, err := json.Marshal(filepaths)
	if err != nil {
		log.Fatal(err)
	}
	errs := notificationSender.Send("watching files "+string(fps), nil)
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			// Non-fatal error when sending notification
		}
	}

	// Create new file watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer func() {
			log.Println("watcher stopped")
			notificationSender.Send("watcher stopped", nil)
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
			log.Fatal(err)
			notificationSender.Send("failed to watch file "+filepath, nil)
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
