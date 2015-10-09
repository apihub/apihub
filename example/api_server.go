package main

import (
	"runtime"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/account/mem"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/db"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	subscription := account.NewEtcdSubscription("/apihub_development", &db.EtcdConfig{Machines: []string{"http://127.0.0.1:2379"}})
	api := api.NewApi(mem.New(), subscription)

	api.AddHook(account.Hook{
		Name:   "apihub-apihub-slack",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: "http://apimanager.org"},
		// Text: `{"username": "ApiHub", "channel": "#apihub",
		// "icon_url": "http://www.albertoleal.me/images/apihub-pq.png",
		// "text": "Um novo serviço foi criado no ApiHub, com o seguinte subdomínio: {{.Service.Subdomain}}."}`,
	})

	api.ListenEvents()
	api.Run()
}
