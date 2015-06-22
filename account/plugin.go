package account

import (
	"github.com/backstage/maestro/errors"
	. "github.com/backstage/maestro/log"
)

type Plugin struct {
	Name    string                 `json:"name"`
	Service string                 `json:"service"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

func (pc *Plugin) Save(service Service) error {
	pc.Service = service.Subdomain

	if err := pc.valid(); err != nil {
		Logger.Info("Failed to save a plugin with invalid data: %+v.", pc)
		return err
	}

	err := store.UpsertPlugin(*pc)
	Logger.Info("plugin.Save: %+v. Err: %s.", pc, err)
	return err
}

func (pc Plugin) Delete() error {
	err := store.DeletePlugin(pc)
	Logger.Info("plugin.Delete: %+v. Err: %s.", pc, err)
	return err
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
