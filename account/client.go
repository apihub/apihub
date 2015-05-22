package account

import (
	"github.com/backstage/backstage/db"
	"github.com/backstage/backstage/errors"
	"github.com/backstage/backstage/util"
	. "github.com/mrvdot/golang-utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The Client type is an encapsulation of a client details.
// It's used but oauth server to grant access token.
type Client struct {
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectUri string `json:"redirect_uri"`
	Owner       string `json:"owner"`
	Team        string `json:"team"`
}

// Save creates a new client.
//
// It requires to inform the fields: Name.
func (client *Client) Save(owner *User, team *Team) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	if client.Name == "" {
		return &errors.ValidationError{Payload: "Name cannot be empty."}
	}
	if client.Id == "" {
		client.Id = GenerateSlug(client.Name)
	} else {
		client.Id = GenerateSlug(client.Id)
	}
	if client.Secret == "" {
		client.Secret = util.GenerateRandomStr(32)
	}

	client.Owner = owner.Email
	client.Team = team.Alias
	es, err := FindClientById(client.Id)
	if err == nil && client.Team == es.Team {
		_, err = conn.Clients().Upsert(bson.M{"id": client.Id}, client)
	} else {
		err = conn.Clients().Insert(client)
	}
	if mgo.IsDup(err) {
		return &errors.ValidationError{Payload: "There is another client with this name."}
	}
	return err
}

// Delete removes an existing client from the server.
func (client *Client) Delete() error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Clients().Remove(client)
	if err == mgo.ErrNotFound {
		return &errors.NotFoundError{Payload: "Client not found."}
	}
	if err != nil {
		return &errors.ValidationError{Payload: err.Error()}
	}
	return err
}

// Try to find a client by its id.
// If the client is not found, return an error. Return the client otherwise.
func FindClientById(id string) (*Client, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var client Client
	err = conn.Clients().Find(bson.M{"id": id}).One(&client)
	if err == mgo.ErrNotFound {
		return nil, &errors.NotFoundError{Payload: "Client not found."}
	}

	return &client, nil
}

// Try to find a client by its id and team.
// If the client is not found, return an error. Return the client otherwise.
func FindClientByIdAndTeam(id, team string) (*Client, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var client Client
	err = conn.Clients().Find(bson.M{"id": id, "team": team}).One(&client)
	if err == mgo.ErrNotFound {
		return nil, &errors.NotFoundError{Payload: "Client not found."}
	}

	return &client, nil
}

// DeleteClientByIdAndTeam removes an existing client from the server based on given id and team.
func DeleteClientByIdAndTeam(id, team string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Clients().Remove(bson.M{"id": id, "team": team})
	if err == mgo.ErrNotFound {
		return &errors.NotFoundError{Payload: "Client not found."}
	}
	return err
}

// DeleteClientByTeam removes an existing client from the server based on given id and team.
func DeleteClientByTeam(team string) error {
	conn, err := db.Conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Clients().RemoveAll(bson.M{"team": team})
	if err != nil {
		return err
	}
	return nil
}

// Find all the clients for a given team alias.
// Return an empty list if nothing is found.
func FindClientsByTeam(team string) ([]*Client, error) {
	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var clients []*Client = []*Client{}
	err = conn.Clients().Find(bson.M{"team": team}).All(&clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}
