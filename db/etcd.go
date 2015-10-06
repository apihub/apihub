package db

import (
	"strings"
	"time"

	"github.com/apihub/apihub/errors"
	. "github.com/apihub/apihub/log"
	"github.com/coreos/go-etcd/etcd"
)

type EtcdConfig struct {
	Machines        []string
	CaFile          string
	CertFile        string
	Consistency     string
	prefixKey       string
	KeyFile         string
	EtcdConsistency string
}

type Etcd struct {
	client    *etcd.Client
	prefixKey string
}

func NewEtcd(prefixKey string, config *EtcdConfig) (*Etcd, error) {
	var (
		client *etcd.Client
		err    error
	)

	if config.CertFile != "" && config.KeyFile != "" {
		if client, err = etcd.NewTLSClient(config.Machines, config.CertFile, config.KeyFile, config.CaFile); err != nil {
			Logger.Error("Failed to connect to Etcd. Error: %+v.", err)
			return nil, err
		}
	} else {
		client = etcd.NewClient(config.Machines)
	}

	// Set the default value if not provided.
	if config.EtcdConsistency == "" {
		config.EtcdConsistency = etcd.STRONG_CONSISTENCY
	}

	if err = client.SetConsistency(config.EtcdConsistency); err != nil {
		Logger.Error("Failed to set Etcd consitency. Error: %+v.", err)
		return nil, err
	}

	return &Etcd{client: client, prefixKey: prefixKey}, nil
}

func (e *Etcd) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *Etcd) SetKey(key string, value string, ttl time.Duration) error {
	k := e.expandKeys(key)
	_, err := e.client.Set(k, value, uint64(ttl/time.Second))
	return handleError(err)
}

// TODO: Need to handle dir!
func (e *Etcd) GetKey(key string) (string, error) {
	k := e.expandKeys(key)
	resp, err := e.client.Get(k, false, false)
	if err != nil {
		return "", handleError(err)
	}

	return resp.Node.Value, nil
}

func (e *Etcd) DeleteKey(key string) error {
	k := e.expandKeys(key)
	_, err := e.client.Delete(k, true)
	return handleError(err)
}

func (e *Etcd) Subscribe(key string, receiverC chan interface{}, doneC chan bool) error {
	k := e.expandKeys(key)
	// To watch for the latest change
	waitIndex := uint64(0)
	for {
		resp, err := e.client.Watch(k, waitIndex, true, nil, doneC)

		if err != nil {
			switch err {
			case etcd.ErrWatchStoppedByUser:
				Logger.Info("Stop watching etcd changes on: %s (%s).", k, err.Error())
				return nil
			default:
				Logger.Error("Failed to connect to etcd: %+v.", err)
				return err
			}
		}
		waitIndex = resp.Node.ModifiedIndex + 1

		select {
		case receiverC <- resp.Node.Value:
		case <-doneC:
			return nil
		}
	}
}

func (e Etcd) expandKeys(keys ...string) string {
	return strings.Join(append([]string{e.prefixKey}, keys...), "/")
}

func handleError(err error) error {
	switch erro := err.(type) {
	case *etcd.EtcdError:
		if erro.ErrorCode == 100 {
			return errors.NewNotFoundError(erro)
		}

		if erro.ErrorCode == 105 {
			return errors.NewDuplicateEntryError(erro)
		}
	}

	return err
}
