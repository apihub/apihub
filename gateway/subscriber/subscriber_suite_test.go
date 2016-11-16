package subscriber_test

import (
	"code.cloudfoundry.org/consuladapter/consulrunner"
	consulapi "github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	consulRunner *consulrunner.ClusterRunner
	consulClient *consulapi.Client
)

var _ = BeforeSuite(func() {
	consulRunner = consulrunner.NewClusterRunner(
		9101+GinkgoParallelNode()*consulrunner.PortOffsetLength,
		1,
		"http",
	)

	consulRunner.Start()
	consulRunner.WaitUntilReady()
})

var _ = AfterSuite(func() {
	consulRunner.Stop()
})

var _ = BeforeEach(func() {
	var err error
	Expect(consulRunner.Reset()).To(Succeed())
	consulClient, err = consulapi.NewClient(&consulapi.Config{Address: string(consulRunner.Address())})
	Expect(err).NotTo(HaveOccurred())
})

func TestPublisher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subscriber Suite")
}
