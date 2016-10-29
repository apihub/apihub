package api_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/albertoleal/requests"
	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
	"github.com/apihub/apihub/apihubfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	var (
		fakeStorage *apihubfakes.FakeStorage
		tmpDir      string
		log         *lagertest.TestLogger

		apihubServer *api.ApihubServer
		server       *httptest.Server
		httpClient   requests.HTTPClient
	)

	BeforeEach(func() {
		var err error
		fakeStorage = new(apihubfakes.FakeStorage)
		log = lagertest.NewTestLogger("apihub-services-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-services-test")
		socketPath := path.Join(tmpDir, fmt.Sprintf("apihub_%d.sock", GinkgoParallelNode()))

		apihubServer = api.New(log, "unix", socketPath, fakeStorage)
		Expect(err).NotTo(HaveOccurred())

		server = httptest.NewServer(apihubServer.Handler())
		httpClient = requests.NewHTTPClient(server.URL)
	})

	AfterEach(func() {
		if tmpDir != "" {
			os.RemoveAll(tmpDir)
		}
	})

	Describe("addService", func() {
		BeforeEach(func() {
			fakeStorage.FindServiceByHandleReturns(apihub.ServiceSpec{}, errors.New("service not found"))
		})

		It("adds a new service", func() {
			headers, code, body, _ := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusCreated,
				Method:         http.MethodPost,
				Path:           "/services",
				Body:           `{"handle":"my-handle", "backends":[{"name":"server-a", "address":"http://server-a"}]}`,
			})

			Expect(stringify(body)).To(Equal(`{"handle":"my-handle","disabled":false,"timeout":0,"backends":[{"name":"server-a","address":"http://server-a","heart_beat_address":"","heart_beat_timeout":0,"heart_beat_retry":0}]}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusCreated))
			Expect(fakeStorage.UpsertServiceCallCount()).To(Equal(1))
		})

		Context("when body is invalid", func() {
			It("returns an error and body is not json", func() {
				headers, code, body, _ := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           "not-a-json",
				})

				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":".*"}`))
				Expect(fakeStorage.UpsertServiceCallCount()).To(Equal(0))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})

			It("returns an error when missing required fields", func() {
				bdy := `{"missing":"handle"}`
				headers, code, body, _ := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           bdy,
				})

				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":"Handle and Backend cannot be empty."}`))
				Expect(fakeStorage.UpsertServiceCallCount()).To(Equal(0))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when checking if there is another service for the same handle", func() {
			var reqArgs requests.Args

			BeforeEach(func() {
				reqArgs = requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           `{"handle":"my-handle", "backends":[{"name":"server-a", "address":"http://server-a"}]}`,
				}
				_, code, _, _ := httpClient.MakeRequest(reqArgs)
				Expect(code).To(Equal(http.StatusCreated))
			})

			It("returns an error when handle is already in use", func() {
				fakeStorage.FindServiceByHandleReturns(apihub.ServiceSpec{Handle: "my-handle"}, nil)
				headers, code, body, _ := httpClient.MakeRequest(reqArgs)

				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":"Handle already in use."}`))
				Expect(fakeStorage.UpsertServiceCallCount()).To(Equal(1))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when storing a service fails", func() {
			BeforeEach(func() {
				fakeStorage.UpsertServiceReturns(errors.New("failed to store service."))
			})

			It("returns an error", func() {
				headers, code, body, _ := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           `{"handle":"my-handle", "backends":[{"name":"server-a", "address":"http://server-a"}]}`,
				})

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":"Failed to add new service."}`))
			})
		})
	})

	Describe("listServices", func() {
		BeforeEach(func() {
			fakeStorage.ServicesReturns([]apihub.ServiceSpec{
				apihub.ServiceSpec{
					Handle: "my-handle",
					Backends: []apihub.BackendInfo{
						apihub.BackendInfo{
							Name:    "server-a",
							Address: "http://server-a",
						},
					},
				},
			}, nil)
		})

		It("lists all existing services", func() {
			headers, code, body, _ := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusOK,
				Method:         http.MethodGet,
				Path:           "/services",
			})

			Expect(stringify(body)).To(Equal(`{"items":[{"handle":"my-handle","disabled":false,"timeout":0,"backends":[{"name":"server-a","address":"http://server-a","heart_beat_address":"","heart_beat_timeout":0,"heart_beat_retry":0}]}],"item_count":1}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusOK))
			Expect(fakeStorage.ServicesCallCount()).To(Equal(1))
		})

		Context("when getting the list of services fails", func() {
			BeforeEach(func() {
				fakeStorage.ServicesReturns(nil, errors.New("failed to list services."))
			})

			It("returns an error", func() {
				headers, code, body, _ := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodGet,
					Path:           "/services",
				})

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":"Failed to retrieve service list."}`))
			})
		})
	})

	Describe("removeService", func() {
		BeforeEach(func() {
			fakeStorage.FindServiceByHandleReturns(apihub.ServiceSpec{Handle: "my-handle"}, nil)
		})

		It("removes a service by handle", func() {
			headers, code, body, _ := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusNoContent,
				Method:         http.MethodDelete,
				Path:           "/services/my-handle",
			})

			Expect(stringify(body)).To(Equal(""))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusNoContent))
			Expect(fakeStorage.RemoveServiceCallCount()).To(Equal(1))
		})

		Context("when removing a service fails", func() {
			BeforeEach(func() {
				fakeStorage.RemoveServiceReturns(errors.New("service not found"))
			})

			It("returns an error", func() {
				headers, code, body, _ := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodDelete,
					Path:           "/services/my-bad-handle",
				})

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(stringify(body)).To(MatchRegexp(`{"error":"bad_request","error_description":"Failed to remove service."}`))
			})
		})
	})
})
