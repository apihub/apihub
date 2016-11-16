package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/gateway"
	"github.com/apihub/apihub/gateway/subscriber"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	port            = flag.String("port", ":8080", "Port to be used")
	consulServerURL = flag.String("consul-server", "http://127.0.0.1:8500", "consul server url")
)

func main() {
	flag.Parse()

	// Configure log
	logger := lager.NewLogger("apihub-gateway")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Configure and start server
	reverseProxyCreator := gateway.NewReverseProxyCreator()
	gw := gateway.New(*port, reverseProxyCreator)

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

	// FIXME: extract this to another method
	// FIXME: another channel to remove services
	// FIXME: handle signal to stop gateway
	// FIXME: add test that shutdown one consul server
	servicesCh := make(chan apihub.ServiceSpec)
	stopCh := make(chan struct{})

	sub := subscriber.NewSubscriber(consulClient)
	go sub.Subscribe(logger, apihub.SERVICES_PREFIX, servicesCh, stopCh)

	go func() {
		logger.Debug("waiting-for-services")
		for {
			select {
			case spec := <-servicesCh:
				var backends []string
				for _, be := range spec.Backends {
					backends = append(backends, be.Address)
				}

				proxySpec := gateway.ReverseProxySpec{
					Handle:      spec.Handle,
					Backends:    backends,
					DialTimeout: time.Duration(spec.Timeout),
				}
				if spec.Disabled {
					gw.RemoveService(logger, spec.Handle)
					logger.Info("service-removed ", lager.Data{"handle": spec.Handle})
				} else {
					gw.AddService(logger, proxySpec)
					logger.Info("service-added", lager.Data{"spec": proxySpec})
				}
			}
		}
	}()

	if err := gw.Start(logger); err != nil {
		panic(fmt.Errorf("Failed to start Apihub Gateway: `%s`.", err))
	}
}
