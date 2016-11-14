package publisher

import (
	"encoding/json"

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

func (p *Publisher) Publish(logger lager.Logger, config apihub.ServiceConfig) error {
	log := logger.Session("publisher")
	log.Debug("start")
	defer log.Debug("end")

	log.Info("publish", lager.Data{"config": config})

	spec, err := json.Marshal(config.ServiceSpec)
	if err != nil {
		log.Error("failed-to-marshal-service-data", err)
		return err
	}

	kvp := &api.KVPair{Key: config.ServiceSpec.Handle, Value: spec}
	_, err = p.client.KV().Put(kvp, nil)
	return err
}
