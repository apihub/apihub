package main

import (
	"github.com/backstage/maestro/account/mem"
	"github.com/backstage/maestro/api"
)

func main() {
	api := api.NewApi(mem.New())
	api.Run()
}
