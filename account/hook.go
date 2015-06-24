package account

import (
	"github.com/backstage/maestro/errors"
	. "github.com/backstage/maestro/log"
	utils "github.com/mrvdot/golang-utils"
)

type HookConfig struct {
	Address string `json:"address"`
	Method  string `json:"method,omitempty"`
}

type Hook struct {
	Name   string     `json:"name"`
	Team   string     `json:"team"`
	Events []string   `json:"events"`
	Config HookConfig `json:"config,omitempty"`
	Text   string     `json:"text,omitempty"`
}

func (w *Hook) Save(team Team) error {
	if err := w.valid(); err != nil {
		Logger.Info("Failed to save hook with invalid data: %+v.", w)
		return err
	}

	w.Name = utils.GenerateSlug(w.Name)
	w.Team = team.Alias

	err := store.UpsertHook(*w)
	Logger.Info("hook.Save: %+v. Err: %s.", w, err)
	return err
}

func (w *Hook) Delete() error {
	err := store.DeleteHook(*w)
	Logger.Info("hook.Delete: %+v. Err: %s.", w, err)
	return err
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
