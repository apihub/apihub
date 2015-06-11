package main

import (
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/api"
)

func main() {
	mem := mem.New()
	store := func() (account.Storable, error) {
		return mem, nil
	}
	api := api.NewApi(store)
	api.Run()
}
