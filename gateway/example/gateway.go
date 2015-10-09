package main

import (
	"runtime"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/db"
	"github.com/apihub/apihub/gateway"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	settings := &gateway.Settings{
		Host: "apimanager.org",
		Port: ":8001",
	}

	one := &account.Service{Endpoint: "http://45.55.32.232:8000", Subdomain: "api", Timeout: 3}
	hw := &account.Service{Endpoint: "http://gohttphelloworld.appspot.com", Subdomain: "helloworld", Timeout: 2}
	services := []*account.Service{one, hw}

	pubsub := account.NewEtcdSubscription("/apihub_development", &db.EtcdConfig{Machines: []string{"http://apihub_etcd_1:2379"}})
	gw := gateway.New(settings, pubsub)
	gw.LoadServices(services)
	gw.RefreshServices()
	gw.Run()
}
