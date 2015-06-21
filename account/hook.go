package account

import (
	"github.com/backstage/maestro/errors"
	utils "github.com/mrvdot/golang-utils"
)

type HookConfig struct {
	URL string `json:"url"`
}

type Hook struct {
	Name   string     `json:"name"`
	Team   string     `json:"team"`
	Events []string   `json:"events"`
	Config HookConfig `json:"config,omitempty"`
}

func (w *Hook) Save(team Team) error {

	if err := w.valid(); err != nil {
		return err
	}

	w.Name = utils.GenerateSlug(w.Name)
	w.Team = team.Alias

	return store.UpsertHook(*w)
}

func (w *Hook) Delete() error {
	return store.DeleteHook(*w)
}

func (w Hook) Exists() bool {
	_, err := FindHookByName(w.Name)
	if err != nil {
		return false
	}
	return true
}

func (w *Hook) valid() error {
	if w.Name == "" && w.Team == "" && len(w.Events) == 0 {
		return errors.NewValidationError(errors.ErrHookMissingRequiredFields)
	}
	return nil
}

func FindHookByName(name string) (*Hook, error) {
	hooks, err := store.FindHookByName(name)
	if err != nil {
		return nil, err
	}
	return &hooks, nil
}

func FindHooksByEvent(event string) ([]Hook, error) {
	hooks, err := store.FindHooksByEvent(event)
	if err != nil {
		return nil, err
	}
	return hooks, nil
}
