package mongore

import (
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	"github.com/fatih/structs"
)

type Mongore struct{}

func New() account.Storable {
	return &Mongore{}
}

func (m *Mongore) SaveToken(k account.TokenKey, expires int, d interface{}) error {
	conn, err := db.Conn()
	if err != nil {
		//FIXME: add log.
		return err
	}
	defer conn.Close()

	conn.Tokens(k.Name, expires, structs.Map(d))
	return nil
}

func (m *Mongore) GetToken(k account.TokenKey, t interface{}) error {
	conn, err := db.Conn()
	if err != nil {
		//FIXME: add log.
		return err
	}
	defer conn.Close()
	err = conn.GetTokenValue(k.Name, t)
	if err != nil {
		return &errors.NotFoundError{Payload: "Token not found."}
	}
	return nil
}

func (m *Mongore) DeleteToken(k account.TokenKey) error {
	conn, err := db.Conn()
	if err != nil {
		//FIXME: add log.
		return err
	}
	defer conn.Close()
	_, err = conn.DeleteToken(k.Name)
	return err
}
