package api

import (
	"github.com/backstage/backstage/db"
	. "gopkg.in/check.v1"
)

func (s *S) TestGetClient(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	err := s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)
	c.Assert(err, IsNil)
	cli, err := s.oAuthStorage.GetClient(osinClient.Id)
	c.Assert(err, IsNil)
	c.Assert(cli, DeepEquals, osinClient)
}

func (s *S) TestSetClient(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	err := s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)
	c.Assert(err, IsNil)
}

func (s *S) TestSaveAuthorize(c *C) {
	defer s.oAuthStorage.RemoveAuthorize(authorizeData.Code)
	err := s.oAuthStorage.SaveAuthorize(authorizeData)
	c.Assert(err, IsNil)
}

func (s *S) TestLoadAuthorize(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	err := s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)
	c.Assert(err, IsNil)
	defer s.oAuthStorage.RemoveAuthorize(authorizeData.Code)
	err = s.oAuthStorage.SaveAuthorize(authorizeData)
	c.Assert(err, IsNil)
	d, err := s.oAuthStorage.LoadAuthorize(authorizeData.Code)
	c.Assert(err, IsNil)
	c.Assert(d, DeepEquals, authorizeData)
}

func (s *S) TestRemoveAuthorize(c *C) {
	err := s.oAuthStorage.SaveAuthorize(authorizeData)
	c.Assert(err, IsNil)
	err = s.oAuthStorage.RemoveAuthorize(authorizeData.Code)
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveAuthorizeWithNonExistingCode(c *C) {
	err := s.oAuthStorage.RemoveAuthorize("non-existing-code")
	c.Assert(err, NotNil)
}

func (s *S) TestSaveAccess(c *C) {
	defer s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	err := s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
}

func (s *S) TestLoadAccess(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	err := s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)
	c.Assert(err, IsNil)
	defer s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	err = s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	t, err := s.oAuthStorage.LoadAccess(accessData.AccessToken)
	c.Assert(err, IsNil)
	accessDataFromOSIN := AccessDataFromOSIN(accessData)
	tFromOSIN := AccessDataFromOSIN(t)
	c.Assert(tFromOSIN, DeepEquals, accessDataFromOSIN)
}

func (s *S) TestLoadAccessWithNonExistingAccessToken(c *C) {
	_, err := s.oAuthStorage.LoadAccess("non-existing-access-token")
	c.Assert(err, NotNil)
}

func (s *S) TestRemoveAccess(c *C) {
	err := s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	err = s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveAccessWithNonExistingAccessToken(c *C) {
	err := s.oAuthStorage.RemoveAccess("non-existing-access-token")
	c.Assert(err, NotNil)
}

func (s *S) TestLoadRefresh(c *C) {
	conn, _ := db.Conn()
	defer conn.Close()
	err := s.oAuthStorage.SetClient(osinClient.Id, osinClient)
	defer conn.Clients().RemoveId(osinClient.Id)
	c.Assert(err, IsNil)
	defer s.oAuthStorage.RemoveAccess(accessData.AccessToken)
	err = s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	t, err := s.oAuthStorage.LoadRefresh(accessData.RefreshToken)
	c.Assert(err, IsNil)
	accessDataFromOSIN := AccessDataFromOSIN(accessData)
	tFromOSIN := AccessDataFromOSIN(t)
	c.Assert(tFromOSIN, DeepEquals, accessDataFromOSIN)
}

func (s *S) TestLoadRefreshTokenWithNonExistingAccessToken(c *C) {
	_, err := s.oAuthStorage.LoadRefresh("non-existing-refresh-token")
	c.Assert(err, NotNil)
}

func (s *S) TestRemoveRefresh(c *C) {
	err := s.oAuthStorage.SaveAccess(accessData)
	c.Assert(err, IsNil)
	err = s.oAuthStorage.RemoveRefresh(accessData.RefreshToken)
	c.Assert(err, IsNil)
}

func (s *S) TestRemoveRefreshTokenWithNonExistingAccessToken(c *C) {
	err := s.oAuthStorage.RemoveRefresh("non-existing-refresh-token")
	c.Assert(err, NotNil)
}
