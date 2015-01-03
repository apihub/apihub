package main

import (
	. "github.com/backstage/backstage/api"
	"github.com/backstage/backstage/log"
)

func main() {
	var api = &Api{
		Config: "config.yaml",
	}
	logger := NewCustomLogger()
	logger.SetLevel(log.DEBUG)
	api.Logger(logger)
	// api.Log().Disable()

	api.Start()
}
