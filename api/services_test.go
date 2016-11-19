package api_test

import (
	"encoding/json"
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
		fakeStorage          *apihubfakes.FakeStorage
		fakeServicePublisher *apihubfakes.FakeServicePublisher
		tmpDir               string
		log                  *lagertest.TestLogger

		apihubServer *api.ApihubServer
		server       *httptest.Server
		httpClient   requests.HTTPClient
	)

	BeforeEach(func() {
		var err error
		fakeStorage = new(apihubfakes.FakeStorage)
		fakeServicePublisher = new(apihubfakes.FakeServicePublisher)
		log = lagertest.NewTestLogger("apihub-services-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-services-test")
		socketPath := path.Join(tmpDir, fmt.Sprintf("apihub_%d.sock", GinkgoParallelNode()))

		apihubServer = api.New(log, "unix", socketPath, fakeStorage, fakeServicePublisher)
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
		It("adds a new service", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusCreated,
				Method:         http.MethodPost,
				Path:           "/services",
				Body:           `{"host":"my-host.apihub.dev", "backends":[{"address":"http://server-a"}]}`,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(ContainSubstring(`{"host":"my-host.apihub.dev","disabled":false,"timeout":0,"backends":[{"address":"http://server-a","disabled":false,"heart_beat_address":"","heart_beat_timeout":0}]}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusCreated))
			Expect(fakeStorage.AddServiceCallCount()).To(Equal(1))
		})

		It("publishes the service", func() {
			spec := apihub.ServiceSpec{
				Host:     "my-host.apihub.dev",
				Disabled: false,
				Backends: []apihub.BackendInfo{
					apihub.BackendInfo{
						Address: "http://server-a",
					},
				},
			}
			body, err := json.Marshal(spec)
			Expect(err).NotTo(HaveOccurred())
			_, _, _, err = httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusCreated,
				Method:         http.MethodPost,
				Path:           "/services",
				Body:           string(body),
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeServicePublisher.PublishCallCount()).To(Equal(1))
			_, prefix, s := fakeServicePublisher.PublishArgsForCall(0)
			Expect(spec).To(Equal(s))
			Expect(prefix).To(Equal(apihub.SERVICES_PREFIX))
		})

		Context("when the service spec is disabled", func() {
			It("does not publish the service", func() {
				spec := apihub.ServiceSpec{
					Host:     "my-host.apihub.dev",
					Disabled: true,
					Backends: []apihub.BackendInfo{
						apihub.BackendInfo{
							Address: "http://server-a",
						},
					},
				}
				body, err := json.Marshal(spec)
				Expect(err).NotTo(HaveOccurred())
				_, _, _, err = httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusCreated,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           string(body),
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeServicePublisher.PublishCallCount()).To(Equal(0))
			})
		})

		Context("when publishing a service fails", func() {
			BeforeEach(func() {
				fakeServicePublisher.PublishReturns(errors.New("failed to publish service"))
			})

			It("returns an error", func() {
				_, _, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           `{"host":"my-host.apihub.dev", "backends":[{"address":"http://server-a"}]}`,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"failed to publish service: 'failed to publish service'"}`))
			})

			It("removes the service from the storage", func() {
				_, _, _, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           `{"host":"my-host.apihub.dev", "backends":[{"address":"http://server-a"}]}`,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeStorage.RemoveServiceCallCount()).To(Equal(1))
				host := fakeStorage.RemoveServiceArgsForCall(0)
				Expect(host).To(Equal("my-host.apihub.dev"))
			})
		})

		Context("when body is invalid", func() {
			It("returns an error and body is not json", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           "not-a-json",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(string(body)).To(MatchRegexp(`{"error":"bad_request","error_description":".*"}`))
				Expect(fakeStorage.AddServiceCallCount()).To(Equal(0))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})

			It("returns an error when missing required fields", func() {
				bdy := `{"missing":"host"}`
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           bdy,
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Host and Backend cannot be empty."}`))
				Expect(fakeStorage.AddServiceCallCount()).To(Equal(0))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when storing a service fails", func() {
			BeforeEach(func() {
				fakeStorage.AddServiceReturns(errors.New("host already in use"))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPost,
					Path:           "/services",
					Body:           `{"host":"my-host.apihub.dev", "backends":[{ "address":"http://server-a"}]}`,
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"failed to add service: 'host already in use'"}`))
			})
		})
	})

	Describe("listServices", func() {
		BeforeEach(func() {
			fakeStorage.ServicesReturns([]apihub.ServiceSpec{
				apihub.ServiceSpec{
					Host: "my-host.apihub.dev",
					Backends: []apihub.BackendInfo{
						apihub.BackendInfo{
							Address: "http://server-a",
						},
					},
				},
			}, nil)
		})

		It("lists all existing services", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusOK,
				Method:         http.MethodGet,
				Path:           "/services",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(ContainSubstring(`{"items":[{"host":"my-host.apihub.dev","disabled":false,"timeout":0,"backends":[{"address":"http://server-a","disabled":false,"heart_beat_address":"","heart_beat_timeout":0}]}],"item_count":1}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusOK))
			Expect(fakeStorage.ServicesCallCount()).To(Equal(1))
		})

		Context("when getting the list of services fails", func() {
			BeforeEach(func() {
				fakeStorage.ServicesReturns(nil, errors.New("failed to list services."))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodGet,
					Path:           "/services",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Failed to retrieve service list."}`))
			})
		})
	})

	Describe("removeService", func() {
		BeforeEach(func() {
			fakeStorage.FindServiceByHostReturns(apihub.ServiceSpec{Host: "my-host.apihub.dev"}, nil)
		})

		It("removes a service by host", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusNoContent,
				Method:         http.MethodDelete,
				Path:           "/services/my-host.apihub.dev",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(BeEmpty())
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusNoContent))
			Expect(fakeStorage.RemoveServiceCallCount()).To(Equal(1))
		})

		It("unpublishes the service", func() {
			_, _, _, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusNoContent,
				Method:         http.MethodDelete,
				Path:           "/services/my-host.apihub.dev",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeServicePublisher.UnpublishCallCount()).To(Equal(1))
			_, prefix, host := fakeServicePublisher.UnpublishArgsForCall(0)
			Expect(prefix).To(Equal(apihub.SERVICES_PREFIX))
			Expect(host).To(Equal("my-host.apihub.dev"))
		})

		Context("when removing a service fails", func() {
			BeforeEach(func() {
				fakeStorage.RemoveServiceReturns(errors.New("service not found"))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodDelete,
					Path:           "/services/my-bad-host",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Failed to remove service."}`))
			})
		})
	})

	Describe("findService", func() {
		BeforeEach(func() {
			fakeStorage.FindServiceByHostReturns(apihub.ServiceSpec{
				Host: "my-host.apihub.dev",
				Backends: []apihub.BackendInfo{
					apihub.BackendInfo{
						Address: "http://server-a",
					},
				},
			}, nil)
		})

		It("finds a service", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusOK,
				Method:         http.MethodGet,
				Path:           "/services/my-host.apihub.dev",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(ContainSubstring(`{"host":"my-host.apihub.dev","disabled":false,"timeout":0,"backends":[{"address":"http://server-a","disabled":false,"heart_beat_address":"","heart_beat_timeout":0}]}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusOK))
			Expect(fakeStorage.FindServiceByHostCallCount()).To(Equal(1))
		})

		Context("when finding a service fails", func() {
			BeforeEach(func() {
				fakeStorage.FindServiceByHostReturns(apihub.ServiceSpec{}, errors.New("failed to find service."))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodGet,
					Path:           "/services/invalid-host",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Failed to find service."}`))
			})
		})
	})

	Describe("updateService", func() {
		BeforeEach(func() {
			fakeStorage.FindServiceByHostReturns(apihub.ServiceSpec{
				Host:     "my-host.apihub.dev",
				Disabled: true,
				Backends: []apihub.BackendInfo{
					apihub.BackendInfo{
						Address: "http://server-a",
					},
				},
			}, nil)
		})

		It("updates a service", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusOK,
				Method:         http.MethodPatch,
				Path:           "/services/my-host.apihub.dev",
				Body:           `{"host":"my-host.apihub.dev", "backends":[{"address":"http://another-server-b"}]}`,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(ContainSubstring(`{"host":"my-host.apihub.dev","disabled":true,"timeout":0,"backends":[{"address":"http://another-server-b","disabled":false,"heart_beat_address":"","heart_beat_timeout":0}]}`))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
			Expect(code).To(Equal(http.StatusOK))
			Expect(fakeStorage.FindServiceByHostCallCount()).To(Equal(1))
			Expect(fakeStorage.UpdateServiceCallCount()).To(Equal(1))
			Expect(fakeServicePublisher.UnpublishCallCount()).To(Equal(1))
			_, prefix, s := fakeServicePublisher.UnpublishArgsForCall(0)
			Expect(prefix).To(Equal(apihub.SERVICES_PREFIX))
			Expect(s).To(Equal("my-host.apihub.dev"))
		})

		Context("when changing the service to be enabled", func() {
			var spec apihub.ServiceSpec

			BeforeEach(func() {
				spec = apihub.ServiceSpec{
					Host:     "my-host.apihub.dev",
					Disabled: false,
					Backends: []apihub.BackendInfo{
						apihub.BackendInfo{
							Address: "http://server-a",
						},
					},
				}

				fakeStorage.FindServiceByHostReturns(spec, nil)
			})

			It("publishes the service", func() {
				_, _, _, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusOK,
					Method:         http.MethodPatch,
					Path:           "/services/my-host.apihub.dev",
					Body:           `{"host":"my-host.apihub.dev", "backends":[{"address":"http://another-server-b"}]}`,
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeServicePublisher.PublishCallCount()).To(Equal(1))
				_, prefix, s := fakeServicePublisher.PublishArgsForCall(0)
				Expect(prefix).To(Equal(apihub.SERVICES_PREFIX))
				Expect(s).To(Equal(spec))
			})
		})

		Context("when finding a service fails", func() {
			BeforeEach(func() {
				fakeStorage.FindServiceByHostReturns(apihub.ServiceSpec{}, errors.New("failed to find service."))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPatch,
					Path:           "/services/my-host.apihub.dev",
					Body:           "{}",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Failed to find service."}`))
			})
		})

		Context("when body is invalid", func() {
			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPatch,
					Path:           "/services/my-host.apihub.dev",
					Body:           "not-a-json",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(string(body)).To(MatchRegexp(`{"error":"bad_request","error_description":".*"}`))
				Expect(fakeStorage.UpdateServiceCallCount()).To(Equal(0))
				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when storing a service fails", func() {
			BeforeEach(func() {
				fakeStorage.UpdateServiceReturns(errors.New("failed to store service."))
			})

			It("returns an error", func() {
				headers, code, body, err := httpClient.MakeRequest(requests.Args{
					AcceptableCode: http.StatusBadRequest,
					Method:         http.MethodPatch,
					Path:           "/services/my-host.apihub.dev",
					Body:           "{}",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(headers["Content-Type"]).To(ContainElement("application/json"))
				Expect(code).To(Equal(http.StatusBadRequest))
				Expect(string(body)).To(ContainSubstring(`{"error":"bad_request","error_description":"Failed to update service."}`))
			})
		})
	})
})
