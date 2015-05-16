package api

import (
	. "github.com/backstage/backstage/account"
	. "gopkg.in/check.v1"
)

func (s *S) TestCollectionSerializer(c *C) {
	teams := []*Team{&Team{Name: "Team", Alias: "Alias", Owner: "alice", Users: []string{}}}

	cs := &CollectionSerializer{
		Items: teams,
		Count: len(teams),
	}
	c.Assert(cs.Serializer(), Equals, `{"items":[{"name":"Team","alias":"Alias","users":[],"owner":"alice"}],"item_count":1}`)
}
