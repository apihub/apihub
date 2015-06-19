package account

import (
	"github.com/backstage/maestro/errors"
)

type Plugin struct {
	Name    string                 `json:"name"`
	Service string                 `json:"service"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func (pc *Plugin) Save(service Service) error {
	pc.Service = service.Subdomain
	if !pc.valid() {
		return errors.NewValidationError(errors.ErrPluginMissingRequiredFields)
	}

	return store.UpsertPlugin(*pc)
}

func FindPluginByNameAndService(pluginName string, service Service) (*Plugin, error) {
	plugin, err := store.FindPluginByNameAndService(pluginName, service)
	if err != nil {
		return nil, err
	}

	return &plugin, nil
}

func (pc Plugin) Delete() error {
	return store.DeletePlugin(pc)
}

func (pc *Plugin) valid() bool {
	if pc.Name == "" || pc.Service == "" {
		return false
	}

	return true
}
