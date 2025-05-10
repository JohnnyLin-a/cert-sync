package apis

import (
	"log"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/johnnylin-a/cert-sync/internal/configs"
)

var notification *router.ServiceRouter

func init() {
	n, err := shoutrrr.CreateSender(configs.GetAppConfig().Notifications...)
	if err != nil {
		panic(err)
	}
	notification = n
}

func GetNotificationsSender() *router.ServiceRouter {
	return notification
}

func LogAndSendNotification(s string) {
	log.Println(s)
	errs := notification.Send(s, nil)
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			// Non-fatal error when sending notification
		}
	}
}
