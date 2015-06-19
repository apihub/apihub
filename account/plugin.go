package account

import "github.com/backstage/maestro/errors"

type Plugin struct {
	Name    string                 `json:"name"`
	Service string                 `json:"service"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func (pc *Plugin) Save(service Service) error {
	pc.Service = service.Subdomain

	if err := pc.valid(); err != nil {
		return err
	}

	return store.UpsertPlugin(*pc)
}

func (pc Plugin) Delete() error {
	return store.DeletePlugin(pc)
}

func (pc *Plugin) valid() error {
	if pc.Name == "" || pc.Service == "" {
		return errors.NewValidationError(errors.ErrPluginMissingRequiredFields)
	}

	return nil
}

func FindPluginByNameAndService(pluginName string, service Service) (*Plugin, error) {
	plugin, err := store.FindPluginByNameAndService(pluginName, service)
	if err != nil {
		return nil, err
	}

	return &plugin, nil
}
