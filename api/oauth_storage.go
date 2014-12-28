package api

import (
	"time"

	"github.com/RangelReale/osin"
	"github.com/backstage/backstage/db"
	"gopkg.in/mgo.v2/bson"
)

type OAuthMongoStorage struct{}

func NewOAuthMongoStorage() *OAuthMongoStorage {
	return &OAuthMongoStorage{}
}

type AuthorizeData struct {
	ClientId    string
	Code        string
	ExpiresIn   int32
	Scope       string
	RedirectUri string
	State       string
	CreatedAt   time.Time
	UserData    interface{}
}

func AuthorizeDataFromOSIN(data *osin.AuthorizeData) *AuthorizeData {
	return &AuthorizeData{
		ClientId:    data.Client.GetId(),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.UserData,
	}
}

func (data *AuthorizeData) AuthorizeDataToOSIN(s *OAuthMongoStorage) (*osin.AuthorizeData, error) {
	client, err := s.GetClient(data.ClientId)
	if err != nil {
		return nil, err
	}
	return &osin.AuthorizeData{
		Client:      client,
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.UserData,
	}, nil
}

type AccessData struct {
	ClientId     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
	Scope        string
	RedirectUri  string
	CreatedAt    time.Time
	UserData     interface{}
}

func AccessDataFromOSIN(data *osin.AccessData) *AccessData {
	return &AccessData{
		ClientId:     data.Client.GetId(),
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}
}

func (data *AccessData) AccessDataToOSIN(s *OAuthMongoStorage) (*osin.AccessData, error) {
	client, err := s.GetClient(data.ClientId)
	if err != nil {
		return nil, err
	}
	return &osin.AccessData{
		Client:       client,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		ExpiresIn:    data.ExpiresIn,
		Scope:        data.Scope,
		RedirectUri:  data.RedirectUri,
		CreatedAt:    data.CreatedAt,
		UserData:     data.UserData,
	}, nil
}

func (s *OAuthMongoStorage) GetClient(id string) (osin.Client, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := new(osin.DefaultClient)
	err = conn.Clients().FindId(id).One(&client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *OAuthMongoStorage) SetClient(id string, client osin.Client) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Clients().UpsertId(id, client)
	return err
}

func (s *OAuthMongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Authorizations().UpsertId(data.Code, AuthorizeDataFromOSIN(data))
	return err
}

func (s *OAuthMongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var authData AuthorizeData
	err = conn.Authorizations().Find(bson.M{"code": code}).One(&authData)
	if err != nil {
		return nil, err
	}
	return authData.AuthorizeDataToOSIN(s)
}

func (s *OAuthMongoStorage) RemoveAuthorize(code string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Authorizations().RemoveId(code)
	return err
}

func (s *OAuthMongoStorage) SaveAccess(data *osin.AccessData) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Accesses().UpsertId(data.AccessToken, AccessDataFromOSIN(data))
	return err
}

func (s *OAuthMongoStorage) LoadAccess(token string) (*osin.AccessData, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var accData AccessData
	err = conn.Accesses().FindId(token).One(&accData)
	if err != nil {
		return nil, err
	}
	return accData.AccessDataToOSIN(s)
}

func (s *OAuthMongoStorage) RemoveAccess(token string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Accesses().RemoveId(token)
	return err
}

func (s *OAuthMongoStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var accData AccessData
	err = conn.Accesses().Find(bson.M{"refreshtoken": token}).One(&accData)
	if err != nil {
		return nil, err
	}
	return accData.AccessDataToOSIN(s)
}

func (s *OAuthMongoStorage) RemoveRefresh(token string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Accesses().Remove(bson.M{"refreshtoken": token})
}

func (s *OAuthMongoStorage) Clone() osin.Storage {
	return s
}

func (s *OAuthMongoStorage) Close() {}
