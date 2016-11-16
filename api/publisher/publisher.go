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

	kvp := &api.KVPair{Key: fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, serviceSpec.Handle), Value: spec}
	_, err = p.client.KV().Put(kvp, nil)
	return err
}

func (p *Publisher) Unpublish(logger lager.Logger, prefix string, handle string) error {
	log := logger.Session("publisher-unpublish")
	log.Debug("start")
	defer log.Debug("end")

	log.Info("unpublish", lager.Data{"handle": handle})

	spec := apihub.ServiceSpec{
		Handle:   handle,
		Disabled: true,
	}
	serviceSpec, err := json.Marshal(spec)
	if err != nil {
		log.Error("failed-to-marshal-service-data", err)
		return err
	}
	kvp := &api.KVPair{Key: fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, spec.Handle), Value: serviceSpec}
	_, err = p.client.KV().Put(kvp, nil)
	if err != nil {
		log.Error("failed-to-unpublish-service", err)
		return err
	}
	key := fmt.Sprintf("%s%s", apihub.SERVICES_PREFIX, handle)
	_, err = p.client.KV().Delete(key, nil)
	return err
}
