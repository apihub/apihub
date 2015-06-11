package api_new_test

import (
	"github.com/backstage/backstage/account_new"
	"github.com/backstage/backstage/api_new"
	. "gopkg.in/check.v1"
)

func (s *S) TestCollectionSerializer(c *C) {
	teams := []*account_new.Team{
		&account_new.Team{Name: "Team", Alias: "Alias", Owner: "alice", Users: []string{}},
	}

	cs := &api_new.CollectionSerializer{
		Items: teams,
		Count: len(teams),
	}
	c.Assert(cs.Serializer(), Equals, `{"items":[{"name":"Team","alias":"Alias","users":[],"owner":"alice"}],"item_count":1}`)
}
