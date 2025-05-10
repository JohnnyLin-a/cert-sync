package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/johnnylin-a/cert-sync/internal/apis"
	"github.com/johnnylin-a/cert-sync/internal/configs"
	"github.com/johnnylin-a/cert-sync/internal/functions"
)

func main() {
	log.Println("Starting cert-sync")

	ps := getFilesToWatch()
	gracefullyStopWatcherChan := functions.WatchFiles(ps, functions.SyncPath)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	gracefullyStopWatcherChan <- struct{}{}
	<-gracefullyStopWatcherChan

	log.Println("Bye")
	apis.GetNotificationsSender().Send("Bye", nil)
}

func getFilesToWatch() []string {
	m := configs.GetAppConfig().SyncFilePaths
	ps := []string{}
	for k := range m {
		ps = append(ps, k)
	}
	return ps
}
