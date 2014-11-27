package main

import (
	"log"

	"github.com/albertoleal/backstage/api"
)

func main() {
	api, err := api.NewApiServer()
	if err != nil {
		log.Fatal("Could not start Backstage API: ", err)
	}
	err = api.RunServer()
	if err != nil {
		log.Fatal("Could not start Backstage API: ", err)
	}
}
