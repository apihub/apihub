package account

import (
  "github.com/backstage/backstage/db"
  "github.com/backstage/backstage/errors"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

// The Client type is an encapsulation of a client details.
// It's used but oauth server to grant access token.
type Client struct {
	Id          string `bson:"_id" json:"id"`
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

  client.Owner = owner.Email
  client.Team = team.Alias
  if err = conn.Clients().Insert(client); err != nil {
    return &errors.ValidationError{Payload: err.Error()}
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
    return &errors.ValidationError{Payload: "Client not found."}
  }
  if err != nil {
    return &errors.ValidationError{Payload: err.Error()}
  }
  return err
}

// Try to find a client by its name and team.
// If the client is not found, return an error. Return the client otherwise.
func FindClientByNameAndTeam(name, team string) (*Client, error) {
  conn, err := db.Conn()
  if err != nil {
    return nil, err
  }
  defer conn.Close()

  var client Client
  err = conn.Clients().Find(bson.M{"name": name, "team": team}).One(&client)
  if err == mgo.ErrNotFound {
    return nil, &errors.ValidationError{Payload: "Client not found."}
  }

  return &client, nil
}

// DeleteClientByNameAndTeam removes an existing client from the server based on given name and team.
func DeleteClientByNameAndTeam(name, team string) error {
  conn, err := db.Conn()
  if err != nil {
    return err
  }
  defer conn.Close()

  err = conn.Clients().Remove(bson.M{"name": name, "team": team})
  if err == mgo.ErrNotFound {
    return &errors.ValidationError{Payload: "Client not found."}
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