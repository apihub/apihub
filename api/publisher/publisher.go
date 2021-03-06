package publisher

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
	"github.com/hashicorp/consul/api"
)

type Publisher struct {
	client *api.Client
}

func NewPublisher(client *api.Client) *Publisher {
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

	kvp := &api.KVPair{Key: fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, serviceSpec.Host), Value: spec}
	_, err = p.client.KV().Put(kvp, nil)
	log.Info("published")
	return err
}

func (p *Publisher) Unpublish(logger lager.Logger, prefix string, host string) error {
	log := logger.Session("publisher-unpublish")
	log.Debug("start")
	defer log.Debug("end")

	log.Info("unpublish", lager.Data{"host": host})

	spec := apihub.ServiceSpec{
		Host:   host,
		Disabled: true,
	}
	serviceSpec, err := json.Marshal(spec)
	if err != nil {
		log.Error("failed-to-marshal-service-data", err)
		return err
	}
	key := fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, spec.Host)
	kvp := &api.KVPair{Key: key, Value: serviceSpec}
	_, err = p.client.KV().Put(kvp, nil)
	if err != nil {
		log.Error("failed-to-unpublish-service", err)
		return err
	}

	_, err = p.client.KV().Delete(key, nil)
	log.Info("unpublished")
	return err
}
