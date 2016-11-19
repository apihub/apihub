package subscriber

import (
	"encoding/json"
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
	"github.com/hashicorp/consul/api"
)

type Subscriber struct {
	client *api.Client
}

func NewSubscriber(client *api.Client) *Subscriber {
	return &Subscriber{
		client: client,
	}
}

func (s *Subscriber) Subscribe(logger lager.Logger, prefix string, servicesCh chan apihub.ServiceSpec, stop <-chan struct{}) error {
	log := logger.Session("subscriber")
	log.Debug("start")
	defer log.Debug("end")

	go func() {
		defer logger.Info("done")

		keys := keySet{}
		queryOpts := &api.QueryOptions{
			WaitIndex: 0,
			WaitTime:  5 * time.Second,
		}

		for {
			select {
			case <-stop:
				return
			default:
			}

			kvPairs, queryMeta, err := s.client.KV().List(prefix, queryOpts)
			if err != nil {
				log.Error("failed-to-retrieve-services", err)
				queryOpts.WaitIndex = 0
				continue
			}

			queryOpts.WaitIndex = queryMeta.LastIndex
			if kvPairs != nil {
				newKeys := newKeySet(kvPairs)
				specs := diff(keys, newKeys)
				for _, spec := range specs {
					servicesCh <- spec
				}

				keys = newKeys
			}
		}
	}()

	select {
	case <-stop:
		logger.Info("stopped")
		close(servicesCh)
	}
	return nil
}

type keySet map[string]*api.KVPair

func newKeySet(pairs api.KVPairs) keySet {
	set := keySet{}
	for _, pair := range pairs {
		set[pair.Key] = pair
	}
	return set
}

func diff(currentSet keySet, newSet keySet) []apihub.ServiceSpec {
	var (
		specs []apihub.ServiceSpec
		spec  apihub.ServiceSpec
	)

	for key, new := range newSet {
		current, ok := currentSet[key]

		// In case nothing has changed
		if ok && current.ModifyIndex == new.ModifyIndex {
			continue
		}

		// In case it's a service to be added/updated
		if !ok || current.ModifyIndex != new.ModifyIndex {
			if err := json.Unmarshal(new.Value, &spec); err != nil {
				continue
			}
			specs = append(specs, spec)
		}
	}

	return specs
}
