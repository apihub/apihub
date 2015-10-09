package main

import (
	"runtime"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/account/mem"
	// "github.com/apihub/apihub/account/mongore"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/db"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// subscription := account.NewEtcdSubscription("/apihub_development", &db.EtcdConfig{Machines: []string{"http://apihub_etcd_1:2379"}})
	subscription := account.NewEtcdSubscription("/apihub_development", &db.EtcdConfig{Machines: []string{"http://127.0.0.1:2379"}})
	// api := api.NewApi(mongore.New(mongore.Config{
	// Host: "apihub_mongo_1:27017",
	// Host:         "127.0.0.1:27017",
	// DatabaseName: "apihub_account_test",
	// }), subscription)
	api := api.NewApi(mem.New(), subscription)

	api.AddHook(account.Hook{
		Name:   "apihub-service-added-slack",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create"},
		Config: account.HookConfig{Address: "https://hooks.slack.com/services/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"},
		// Config: account.HookConfig{Address: "http://localhost:9999"},
		Text: `{"username": "ApiHub", "channel": "#general",
		"icon_emoji": ":ghost:",
		"text": "A new service has been added with the following subdomain: {{.Service.Subdomain}}."}`,
	})

	api.AddHook(account.Hook{
		Name:   "apihub-service-deleted-slack",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.delete"},
		Config: account.HookConfig{Address: "https://hooks.slack.com/services/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"},
		Text: `{"username": "ApiHub", "channel": "#general",
		"icon_emoji": ":ghost:",
		"text": "Service -> {{.Service.Subdomain}} <- has been removed ."}`,
	})

	api.ListenEvents()
	api.Run()
}
