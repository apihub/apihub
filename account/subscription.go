package account

import (
	"fmt"

	"github.com/apihub/apihub/db"
	. "github.com/apihub/apihub/log"
)

type PubSub interface {
	Publish(name string, data []byte) error
	Subscribe(name string, receiver chan interface{}, done chan bool)
}

var pubsub PubSub

func NewPubSub(ps PubSub) {
	pubsub = ps
}

type EtcdSubscription struct {
	client *db.Etcd
}

func NewEtcdSubscription(prefixKey string, config *db.EtcdConfig) PubSub {
	cli, err := db.NewEtcd(prefixKey, config)
	if err != nil {
		Logger.Error("Failed to establish a connection with Etcd: %+v.", err)
	}

	return &EtcdSubscription{client: cli}
}

func (e *EtcdSubscription) Publish(name string, data []byte) error {
	return e.client.SetKey(name, string(data), 0)
}

func (e *EtcdSubscription) Subscribe(name string, receiverC chan interface{}, done chan bool) {
	go func() {
		if err := e.client.Subscribe(name, receiverC, done); err != nil {
			fmt.Printf("%+v", err)
			close(receiverC)
			close(done)
		}
	}()
}
