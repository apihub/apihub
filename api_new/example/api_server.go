package main

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/account_new/mem"
	"github.com/backstage/backstage/api_new"
)

func main() {
	mem := mem.New()
	store := func() (account_new.Storable, error) {
		return mem, nil
	}
	api := api_new.NewApi(store)
	api.Run()
}
