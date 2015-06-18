package account_test

import (
	"github.com/backstage/maestro/account"
	"github.com/backstage/maestro/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestSaveWebhook(c *C) {
	err := webhook.Save(team)
	c.Assert(err, IsNil)
	defer webhook.Delete()
}

func (s *S) TestSaveWebhookWithoutRequiredFields(c *C) {
	webhook = account.Webhook{}
	err := webhook.Save(team)
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestDeleteWebhook(c *C) {
	webhook.Save(team)
	c.Assert(webhook.Exists(), Equals, true)
	webhook.Delete()
	c.Assert(webhook.Exists(), Equals, false)
}

func (s *S) TestFindWebhookByName(c *C) {
	webhook.Save(team)

	t, err := account.FindWebhookByName(webhook.Name)
	c.Check(t, Not(IsNil))
	c.Check(err, IsNil)
	defer webhook.Delete()
}

func (s *S) TestFindWebhookByNameNotFound(c *C) {
	t, err := account.FindWebhookByName("not-found")
	c.Check(t, IsNil)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}
