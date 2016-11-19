package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/api/publisher"
	"github.com/apihub/apihub/storage"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	network         = flag.String("network", "unix", "Either `tcp` or `unix`")
	address         = flag.String("address", "/tmp/apihub.sock", "Port for `tcp` or filepath for `unix`")
	consulServerURL = flag.String("consul-server", "http://127.0.0.1:8500", "consul server url")
)

func main() {
	flag.Parse()

	// Configure log
	logger := lager.NewLogger("apihub-api")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Configure and start server
	store := storage.New()
	consulURL, err := url.Parse(*consulServerURL)
	if err != nil {
		panic(fmt.Sprintf("Error parsing Consul URL: %s", err))
	}
	consulClient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulURL.Host,
		Scheme:  consulURL.Scheme,
	})
	if err != nil {
		panic(fmt.Sprintf("Error connecting to Consul agent: %s", err))
	}
	publisher := publisher.NewPublisher(consulClient)
	server := api.New(logger, *network, *address, store, publisher)
	server.Start(true)
}
