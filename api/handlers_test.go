package api_test

import (
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
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When a client connects", func() {
	var (
		fakeStorage          *apihubfakes.FakeStorage
		fakeServicePublisher *apihubfakes.FakeServicePublisher
		apihubClient         apihub.Client
		apihubServer         *api.ApihubServer
		err                  error
		log                  *lagertest.TestLogger
		tmpDir               string
	)

	BeforeEach(func() {
		fakeStorage = new(apihubfakes.FakeStorage)
		fakeServicePublisher = new(apihubfakes.FakeServicePublisher)
		log = lagertest.NewTestLogger("apihub-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-test")
		socketPath := path.Join(tmpDir, "apihub.sock")
		apihubServer = api.New(log, "unix", socketPath, fakeStorage, fakeServicePublisher)
		Expect(err).NotTo(HaveOccurred())

		Expect(apihubServer.Start(false)).NotTo(HaveOccurred())
		apihubClient = client.New(connection.New("unix", socketPath))
	})

	AfterEach(func() {
		if tmpDir != "" {
			os.RemoveAll(tmpDir)
		}
	})

	Describe("and sends a Ping Request", func() {
		Context("and the server is up and running", func() {
			It("does not return an error", func() {
				Expect(apihubClient.Ping()).NotTo(HaveOccurred())
			})
		})

		Context("when the server is not up and running", func() {
			BeforeEach(func() {
				Expect(apihubServer.Stop()).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				Expect(apihubClient.Ping()).To(HaveOccurred())
			})
		})
	})
})

var _ = Describe("When an http request is sent", func() {
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
		log = lagertest.NewTestLogger("apihub-handler-test")
		tmpDir, err = ioutil.TempDir(os.TempDir(), "apihub-server-handler-test")
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

	Context("when the page is not found", func() {
		It("returns page not found", func() {
			headers, code, body, err := httpClient.MakeRequest(requests.Args{
				AcceptableCode: http.StatusNotFound,
				Method:         http.MethodGet,
				Path:           "/not-found",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(ContainSubstring(`{"error":"not_found","error_description":"The resource does not exist."}`))
			Expect(code).To(Equal(http.StatusNotFound))
			Expect(headers["Content-Type"]).To(ContainElement("application/json"))
		})
	})
})
