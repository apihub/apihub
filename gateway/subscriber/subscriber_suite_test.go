package subscriber_test

import (
	"code.cloudfoundry.org/consuladapter"
	"code.cloudfoundry.org/consuladapter/consulrunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	consulRunner *consulrunner.ClusterRunner
	consulClient consuladapter.Client
)

func TestPublisher(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		consulRunner = consulrunner.NewClusterRunner(
			9101+GinkgoParallelNode()*consulrunner.PortOffsetLength,
			1,
			"http",
		)

		consulRunner.Start()
		consulRunner.WaitUntilReady()
		return nil
	}, func(_ []byte) {
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		consulRunner.Stop()
	})

	RunSpecs(t, "Subscriber Suite")
}

var _ = BeforeEach(func() {
	consulClient = consulRunner.NewClient()
})

var _ = AfterEach(func() {
	Expect(consulRunner.Reset()).To(Succeed())
})
