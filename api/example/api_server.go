package main

import (
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/api"
)

func main() {
	api := api.NewApi(mem.New())
	api.Run()
}
