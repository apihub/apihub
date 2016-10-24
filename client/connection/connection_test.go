package connection_test

import (
	"net/http"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/client/connection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Connection", func() {

	var (
		conn    connection.Connection
		server  *ghttp.Server
		network string
		address string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		network = "tcp"
		address = server.Addr()
	})

	JustBeforeEach(func() {
		conn = connection.New(network, address)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Ping", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/ping"),
						ghttp.RespondWith(200, `{"ping":"pong"}`),
					),
				)
			})

			It("pings the service", func() {
				err := conn.Ping()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, api.Routes[api.Ping].Path),
						ghttp.RespondWith(500, "{}"),
					),
				)
			})

			It("returns an error", func() {
				Expect(conn.Ping()).To(HaveOccurred())
			})
		})
	})

	Describe("AddService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, api.Routes[api.AddService].Path),
						ghttp.RespondWith(201, `{"handle":"my-handle"}`),
					),
				)
			})

			It("adds a new service", func() {
				spec, err := conn.AddService(
					apihub.ServiceSpec{
						Handle: "my-handle",
					},
				)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Handle).To(Equal("my-handle"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, api.Routes[api.AddService].Path),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.AddService(apihub.ServiceSpec{})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
