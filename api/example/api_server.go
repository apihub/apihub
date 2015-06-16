package main

import (
	"github.com/backstage/apimanager/account/mem"
	"github.com/backstage/apimanager/api"
)

func main() {
	api := api.NewApi(mem.New())
	api.Run()
}
