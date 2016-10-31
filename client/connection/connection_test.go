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

	Describe("Services", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, api.Routes[api.ListServices].Path),
						ghttp.RespondWith(200, `{"items":[{"handle":"my-handle"}, {"handle":"another-handle"}],"item_count":1}`),
					),
				)
			})

			It("lists existing services specs", func() {
				specs, err := conn.Services()

				Expect(err).NotTo(HaveOccurred())
				Expect(len(specs)).To(Equal(2))
				Expect(specs).To(ConsistOf([]apihub.ServiceSpec{
					apihub.ServiceSpec{Handle: "my-handle"},
					apihub.ServiceSpec{Handle: "another-handle"},
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
						ghttp.VerifyRequest(http.MethodDelete, "/services/my-handle"),
						ghttp.RespondWith(204, ""),
					),
				)
			})

			It("removes an existing service", func() {
				err := conn.RemoveService("my-handle")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/services/my-handle"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				err := conn.RemoveService("my-handle")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/services/my-handle"),
						ghttp.RespondWith(200, `{"handle":"my-handle"}`),
					),
				)
			})

			It("finds a service", func() {
				spec, err := conn.FindService("my-handle")

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Handle).To(Equal("my-handle"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/services/invalid-handle"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.FindService("invalid-handle")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UpdateService", func() {
		Context("when the request succeeds", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/services/my-handle"),
						ghttp.RespondWith(200, `{"handle":"my-handle","disabled":true}`),
					),
				)
			})

			It("updates the service", func() {
				spec := apihub.ServiceSpec{
					Handle:   "my-handle",
					Disabled: true,
				}

				spec, err := conn.UpdateService("my-handle", spec)

				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Handle).To(Equal("my-handle"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/services/invalid-handle"),
						ghttp.RespondWith(400, "{}"),
					),
				)
			})

			It("returns an error", func() {
				_, err := conn.UpdateService("invalid-handle", apihub.ServiceSpec{})
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
