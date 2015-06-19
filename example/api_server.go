package main

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/api"
)

func main() {
	api := api.NewApi(mem.New())
	api.AddWebhook(account.Webhook{
		Name:   "maestro-gateway-services",
		Team:   "*",
		Events: []string{"service.create", "service.update", "service.delete"},
		Config: account.WebhookConfig{Url: "http://localhost:8001"},
	})

	api.Run()
}
