package main

import (
	"runtime"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/api"
	"github.com/backstage/maestro/db"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	subscription := account.NewEtcdSubscription("/maestro_development", &db.EtcdConfig{Machines: []string{"http://localhost:2379"}})
	api := api.NewApi(mem.New(), subscription)

	api.AddHook(account.Hook{
		Name:   "backstage-maestro-slack",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: "http://localhost:8001"},
		// Text: `{"username": "Backstage Maestro", "channel": "#backstage",
		// "icon_url": "http://www.albertoleal.me/images/maestro-pq.png",
		// "text": "Um novo serviço foi criado no Backstage Maestro, com o seguinte subdomínio: {{.Service.Subdomain}}."}`,
	})

	api.ListenEvents()
	api.Run()
}
