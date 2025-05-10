package functions

import (
	"log"

	"github.com/johnnylin-a/cert-sync/internal/configs"
)

func SyncPath(path string) {
	log.Println("Sync " + path)
	cfg := configs.GetAppConfig()
	if syncCfg, ok := cfg.SyncFilePaths[path]; ok {
		for _, v := range syncCfg {
			log.Println("Will copy " + path + " on ssh host " + v.ConfigName + " to path " + v.Dst)
		}
	}
}
