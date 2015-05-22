package db

import . "gopkg.in/check.v1"

func (s *S) TestServices(c *C) {
	storage, err := Conn()
	c.Assert(err, IsNil)
	services := storage.Services()
	collection := storage.Collection("services")
	c.Assert(services, DeepEquals, collection)
}

func (s *S) TestUsers(c *C) {
	storage, err := Conn()
	c.Assert(err, IsNil)
	users := storage.Users()
	collection := storage.Collection("users")
	c.Assert(users, DeepEquals, collection)
}

func (s *S) TestTeams(c *C) {
	storage, err := Conn()
	c.Assert(err, IsNil)
	teams := storage.Teams()
	collection := storage.Collection("teams")
	c.Assert(teams, DeepEquals, collection)
}

func (s *S) TestPluginsConfig(c *C) {
	storage, err := Conn()
	c.Assert(err, IsNil)
	midds := storage.PluginsConfig()
	collection := storage.Collection("plugins_config")
	c.Assert(midds, DeepEquals, collection)
}
