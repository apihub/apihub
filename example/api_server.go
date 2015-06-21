package main

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/api"
)

func main() {

	api := api.NewApi(mem.New())

	api.AddHook(account.Hook{
		Name:   "maestro-gateway-services",
		Team:   account.ALL_TEAMS,
		Events: []string{"service.create", "service.update", "service.delete"},
		Config: account.HookConfig{URL: "http://localhost:8001"},
	})

	api.ListenEvents()
	api.Run()
}
