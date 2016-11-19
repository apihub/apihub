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

	Context("when the request fails", func() {
		It("returns the default message when failing to parse response", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, api.Routes[api.Ping].Path),
					ghttp.RespondWith(400, `this is not a valid json`),
				),
			)

			err := conn.Ping()
			Expect(err).To(MatchError(ContainSubstring("request failed")))
		})

		It("returns the error message when succeding to parse response", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, api.Routes[api.Ping].Path),
					ghttp.RespondWith(400, `{"error":"bad_request", "error_description":"Failed to ping."}`),
				),
			)

			err := conn.Ping()
			Expect(err).To(MatchError(ContainSubstring("Failed to ping.")))
		})
	})

	Describe("Ping", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, api.Routes[api.Ping].Path),
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
						ghttp.RespondWith(201, `{"host":"my-host"}`),
					),
				)
			})

			It("adds a new service", func() {
				spec, err := conn.AddService(
					apihub.ServiceSpec{
						Host: "my-host",
					},
				)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Host).To(Equal("my-host"))
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

	Describe("Services", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, api.Routes[api.ListServices].Path),
						ghttp.RespondWith(200, `{"items":[{"host":"my-host"}, {"host":"another-host"}],"item_count":1}`),
					),
				)
			})

			It("lists existing services specs", func() {
				specs, err := conn.Services()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(specs)).To(Equal(2))
				Expect(specs).To(ConsistOf([]apihub.ServiceSpec{
					apihub.ServiceSpec{Host: "my-host"},
					apihub.ServiceSpec{Host: "another-host"},
				}))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, api.Routes[api.ListServices].Path),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.Services()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("RemoveService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/services/my-host"),
						ghttp.RespondWith(204, ""),
					),
				)
			})

			It("removes an existing service", func() {
				err := conn.RemoveService("my-host")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/services/my-host"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				err := conn.RemoveService("my-host")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/services/my-host"),
						ghttp.RespondWith(200, `{"host":"my-host"}`),
					),
				)
			})

			It("finds a service", func() {
				spec, err := conn.FindService("my-host")

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Host).To(Equal("my-host"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/services/invalid-host"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.FindService("invalid-host")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UpdateService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/services/my-host"),
						ghttp.RespondWith(200, `{"host":"my-host","disabled":true}`),
					),
				)
			})

			It("updates the service", func() {
				spec := apihub.ServiceSpec{
					Host:     "my-host",
					Disabled: true,
				}

				spec, err := conn.UpdateService("my-host", spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Host).To(Equal("my-host"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/services/invalid-host"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.UpdateService("invalid-host", apihub.ServiceSpec{})
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
