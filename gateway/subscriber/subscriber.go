package subscriber

import (
	"encoding/json"
	"time"

	"code.cloudfoundry.org/consuladapter"
	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
	"github.com/hashicorp/consul/api"
)

type Subscriber struct {
	client consuladapter.Client
}

func NewSubscriber(client consuladapter.Client) *Subscriber {
	return &Subscriber{
		client: client,
	}
}

func (s *Subscriber) Subscribe(logger lager.Logger, prefix string, servicesCh chan apihub.ServiceSpec, stop <-chan struct{}) error {
	log := logger.Session("subscriber")
	log.Debug("start")
	defer log.Debug("end")

	go func() {
		defer logger.Info("finished")

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
				diffKeys := diff(keys, newKeys)
				for _, kv := range diffKeys {
					var spec apihub.ServiceSpec
					if err := json.Unmarshal(kv.Value, &spec); err != nil {
						log.Error("failed-to-parse-spec", err, lager.Data{"key": kv.Key})
						continue
					}
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

func diff(currentSet keySet, newSet keySet) []*api.KVPair {
	var missing []*api.KVPair
	for key, new := range newSet {
		current, ok := currentSet[key]
		if !ok || current.ModifyIndex != new.ModifyIndex {
			missing = append(missing, new)
		}
	}

	return missing
}
