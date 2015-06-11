package api_test

import (
	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/api"
	. "gopkg.in/check.v1"
)

func (s *S) TestCollectionSerializer(c *C) {
	teams := []*account.Team{
		&account.Team{Name: "Team", Alias: "Alias", Owner: "alice", Users: []string{}},
	}

	cs := &api.CollectionSerializer{
		Items: teams,
		Count: len(teams),
	}
	c.Assert(cs.Serializer(), Equals, `{"items":[{"name":"Team","alias":"Alias","users":[],"owner":"alice"}],"item_count":1}`)
}
