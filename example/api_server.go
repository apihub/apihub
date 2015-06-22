package main

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/api"
)

func main() {

	api := api.NewApi(mem.New())

	// api.AddHook(account.Hook{
	// 	Name:   "maestro-gateway-services",
	// 	Team:   account.ALL_TEAMS,
	// 	Events: []string{"service.create", "service.update", "service.delete"},
	// 	Config: account.HookConfig{Address: "http://localhost:8001"},
	// })

	api.AddHook(account.Hook{
		Name:   "backstage-slack-bla",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create", "service.update", "service.delete"},
		Config: account.HookConfig{Address: "http://localhost:8001"},
		Text: `{"username": "Backstage Maestro", "channel": "#backstage",
		"icon_url": "http://www.albertoleal.me/images/maestro-pq.png",
		"text": "Um novo serviço foi criado no Backstage Maestro, com o seguinte subdomínio: {{.Service.Subdomain}}."}`,
	})

	api.AddHook(account.Hook{
		Name:   "backstage-slack-services",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create", "service.update", "service.delete"},
		Config: account.HookConfig{Address: "http://localhost:8002"},
		Text:   `{"subdomain": "{{.Service.Subdomain}}", "endpoint": "{{.Service.Endpoint}}","disabled": "{{.Service.Disabled}}"}`,
	})

	api.ListenEvents()
	api.Run()
}
