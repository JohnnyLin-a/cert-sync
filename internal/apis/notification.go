package apis

import (
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
