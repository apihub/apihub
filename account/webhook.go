package account

import (
	"github.com/backstage/maestro/errors"
	utils "github.com/mrvdot/golang-utils"
)

type WebhookConfig struct {
	Url string `json:"url"`
}

type Webhook struct {
	Name   string        `json:"name"`
	Team   string        `json:"team"`
	Events []string      `json:"events"`
	Config WebhookConfig `json:"config,omitempty"`
}

func (w *Webhook) Save(team Team) error {
	if !w.valid() {
		return errors.NewValidationError(errors.ErrWebhookMissingRequiredFields)
	}
	w.Name = utils.GenerateSlug(w.Name)
	w.Team = team.Alias

	return store.UpsertWebhook(*w)
}

func (w *Webhook) Delete() error {
	return store.DeleteWebhook(*w)
}

func (w Webhook) Exists() bool {
	_, err := store.FindWebhookByName(w.Name)
	if err != nil {
		return false
	}
	return true
}

func (w *Webhook) valid() bool {
	if w.Name == "" && w.Team == "" && len(w.Events) == 0 {
		return false
	}
	return true
}

func FindWebhookByName(name string) (*Webhook, error) {
	webhook, err := store.FindWebhookByName(name)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}
