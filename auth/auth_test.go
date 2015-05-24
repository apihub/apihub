package auth

import (
	"testing"
	"time"

	"github.com/backstage/backstage/account"
	"github.com/backstage/backstage/account/mem"
	"github.com/backstage/backstage/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = Describe("Token", func() {
	var authToken AuthenticationToken
	var user account.User
	var token account.TokenInfo
	var t *account.TokenInfo
	BeforeEach(func() {
		authToken = AuthenticationToken{
			storage: mem.New(),
		}

		user = account.User{Name: "Alice", Username: "alice", Email: "alice@example.org", Password: "123456"}
	})

	Describe("GetUserFromToken", func() {
		BeforeEach(func() {
			token = account.TokenInfo{
				Token:     "valid-token",
				Expires:   10,
				CreatedAt: time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z07:00"),
				Type:      "Token",
			}

			err := authToken.storage.SaveToken(account.TokenKey{Name: token.Token}, &user)
			if err != nil {
				panic(err)
			}
		})

		It("returns an error when the token is invalid", func() {
			user, err := authToken.GetUserFromToken("Bearer invalid")
			Expect(user).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.ErrInvalidTokenFormat))
		})

		It("returns an error when the token is not found", func() {
			user, err := authToken.GetUserFromToken("Token not-found")
			Expect(user).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(&errors.NotFoundError{Payload: "Token not found."}))
		})

		It("returns an instance of user when the token is found", func() {
			user, err := authToken.GetUserFromToken("Token valid-token")
			Expect(user).To(Equal(user))
			Expect(err).To(BeNil())
		})
	})

	Describe("TokenFor", func() {
		It("generates a token for given user", func() {
			t, err := authToken.TokenFor(&user)
			Expect(t).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
			Expect(len(t.Token)).To(Equal(44))
		})

		It("returns the same token, if valid, instead of generating another", func() {
			t, err := authToken.TokenFor(&user)
			Expect(t).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())

			t2, err := authToken.TokenFor(&user)
			Expect(t.Token).To(Equal(t2.Token))
		})
	})

	Describe("RevokeTokenFor", func() {
		BeforeEach(func() {
			t, _ = authToken.TokenFor(&user)
		})

		It("revokes current token for given user", func() {
			tk := t.Type + " " + t.Token
			u, err := authToken.GetUserFromToken(tk)
			Expect(u).NotTo(BeNil())

			authToken.RevokeTokenFor(&user)
			u, err = authToken.GetUserFromToken(tk)
			Expect(u).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})
})
