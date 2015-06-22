package account

import (
	"github.com/backstage/maestro/errors"
	. "github.com/backstage/maestro/log"
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

	if err := app.valid(); err != nil {
		Logger.Info("Failed to create an app with invalid data: %+v.", app)
		return err
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
		Logger.Info("Failed to create an app with duplicate data: %+v.", app)
		return errors.NewValidationError(errors.ErrAppDuplicateEntry)
	}

	err := store.UpsertApp(*app)
	Logger.Info("app.Create: %+v. Err: %s.", app, err)
	return err
}

func (app *App) Update() error {
	if err := app.valid(); err != nil {
		Logger.Info("Failed to create an app with invalid data: %+v.", app)
		return err
	}

	if !app.Exists() {
		Logger.Info("Failed to update a not-found app: %+v.", app)
		return errors.NewNotFoundError(errors.ErrAppNotFound)
	}

	err := store.UpsertApp(*app)
	Logger.Info("app.Update: %+v. Err: %s.", app, err)
	return err
}

func (app App) Delete(owner User) error {
	if app.Owner != owner.Email {
		Logger.Info("Failed to delete an app. Only the owner has permission to do that: %+v.", app)
		return errors.NewForbiddenError(errors.ErrOnlyOwnerHasPermission)
	}

	err := store.DeleteApp(app)
	Logger.Info("app.Delete: %+v. Err: %s.", app, err)
	return err
}

func (app App) Exists() bool {
	_, err := FindAppByClientId(app.ClientId)
	if err != nil {
		return false
	}
	return true
}

func (app *App) valid() error {
	if app.Name == "" {
		return errors.NewValidationError(errors.ErrAppMissingRequiredFields)
	}

	return nil
}

func DeleteAppsByTeam(team Team, owner User) error {
	apps, err := store.TeamApps(team)
	Logger.Debug("Apps to be deleted: %+v.", apps)
	if err != nil {
		return err
	}
	for _, s := range apps {
		s.Delete(owner)
	}
	Logger.Info("All apps were excluded from the team `%s`.", team.Alias)
	return nil
}

func FindAppByClientId(clientId string) (*App, error) {
	app, err := store.FindAppByClientId(clientId)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
