package main

import (
	"runtime"

	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/db"
	"github.com/backstage/maestro/gateway"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	settings := &gateway.Settings{
		Host: "backstage.example.org",
		Port: ":8001",
	}

	one := &account.Service{Endpoint: "http://localhost:9999", Subdomain: "one", Timeout: 2}
	hw := &account.Service{Endpoint: "http://gohttphelloworld.appspot.com", Subdomain: "helloworld", Timeout: 2}
	services := []*account.Service{one, hw}

	pubsub := account.NewEtcdSubscription("/maestro_development", &db.EtcdConfig{Machines: []string{"http://localhost:2379"}})
	gw := gateway.New(settings, pubsub)
	gw.LoadServices(services)
	gw.RefreshServices()
	gw.Run()
}
