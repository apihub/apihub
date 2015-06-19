package account

import (
	"github.com/backstage/maestro/errors"
	"github.com/backstage/maestro/util"
	goutils "github.com/mrvdot/golang-utils"
)

type App struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Name         string   `json:"name"`
	RedirectUris []string `json:"redirect_uris"`
	Owner        string   `json:"owner"`
	Team         string   `json:"team"`
}

func (app *App) Create(owner User, team Team) error {
	app.Owner = owner.Email
	app.Team = team.Alias

	if !app.valid() {
		return errors.NewValidationError(errors.ErrAppMissingRequiredFields)
	}

	if app.ClientId == "" {
		app.ClientId = goutils.GenerateSlug(app.Name)
	} else {
		app.ClientId = goutils.GenerateSlug(app.ClientId)
	}
	if app.ClientSecret == "" {
		app.ClientSecret = util.GenerateRandomStr(32)
	}

	if app.Exists() {
		return errors.NewValidationError(errors.ErrAppDuplicateEntry)
	}

	return store.UpsertApp(*app)
}

func (app *App) Update() error {
	if !app.valid() {
		return errors.NewValidationError(errors.ErrAppMissingRequiredFields)
	}

	if !app.Exists() {
		return errors.NewNotFoundError(errors.ErrAppNotFound)
	}

	return store.UpsertApp(*app)
}

func (app App) Delete(owner User) error {
	if app.Owner != owner.Email {
		return errors.NewForbiddenError(errors.ErrOnlyOwnerHasPermission)
	}

	err := store.DeleteApp(app)

	return err
}

func DeleteAppsByTeam(team Team, owner User) error {
	apps, err := store.TeamApps(team)
	if err != nil {
		return err
	}
	for _, s := range apps {
		s.Delete(owner)
	}
	return nil
}

func (app App) Exists() bool {
	_, err := store.FindAppByClientId(app.ClientId)
	if err != nil {
		return false
	}
	return true
}

func FindAppByClientId(clientId string) (*App, error) {
	app, err := store.FindAppByClientId(clientId)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (app *App) valid() bool {
	if app.Name == "" {
		return false
	}
	return true
}
