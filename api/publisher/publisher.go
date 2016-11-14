package publisher

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/consuladapter"
	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
	"github.com/hashicorp/consul/api"
)

type Publisher struct {
	client consuladapter.Client
}

func NewPublisher(client consuladapter.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) Publish(logger lager.Logger, prefix string, serviceSpec apihub.ServiceSpec) error {
	log := logger.Session("publisher-publish")
	log.Debug("start")
	defer log.Debug("end")

	log.Info("publish", lager.Data{"serviceSpec": serviceSpec})

	spec, err := json.Marshal(serviceSpec)
	if err != nil {
		log.Error("failed-to-marshal-service-data", err)
		return err
	}

	kvp := &api.KVPair{Key: fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, serviceSpec.Handle), Value: spec}
	_, err = p.client.KV().Put(kvp, nil)
	return err
}
