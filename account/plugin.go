package account

import "github.com/backstage/maestro/errors"

type PluginConfig struct {
	Name    string                 `json:"name"`
	Service string                 `json:"service"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func (pc *PluginConfig) Save(service Service) error {
	pc.Service = service.Subdomain

	if err := pc.valid(); err != nil {
		return err
	}

	return store.UpsertPluginConfig(*pc)
}

func FindPluginByNameAndService(pluginName string, service Service) (*PluginConfig, error) {
	plugin, err := store.FindPluginConfigByNameAndService(pluginName, service)
	if err != nil {
		return nil, err
	}

	return &plugin, nil
}

func (pc PluginConfig) Delete() error {
	return store.DeletePluginConfig(pc)
}

func (pc *PluginConfig) valid() error {
	if pc.Name == "" || pc.Service == "" {
		return errors.NewValidationError(errors.ErrPluginConfigMissingRequiredFields)
	}

	return nil
}
