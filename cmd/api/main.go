package main

import (
	"flag"
	"fmt"
	"os"

	"code.cloudfoundry.org/consuladapter"
	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/publisher"
	"github.com/apihub/apihub/storage"
)

var (
	network         = flag.String("network", "unix", "Either `tcp` or `unix`")
	address         = flag.String("address", "/tmp/apihub.sock", "Port for `tcp` or filepath for `unix`")
	consulServerURL = flag.String("consul-server", "http://127.0.0.1:9999", "consul server url")
)

func main() {
	flag.Parse()

	// Configure log
	logger := lager.NewLogger("apihub-api")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Configure and start server
	store := storage.New()
	consulClient, err := consuladapter.NewClientFromUrl(*consulServerURL)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to Consul agent: %s", err))
	}
	publisher := publisher.NewPublisher(consulClient)
	server := api.New(logger, *network, *address, store, publisher)
	server.Start(true)
}
