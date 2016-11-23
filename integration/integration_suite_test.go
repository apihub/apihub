package integration_test

import (
	"time"

	"code.cloudfoundry.org/consuladapter/consulrunner"

	"github.com/apihub/apihub/integration/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	ApihubAPIBin     string
	ApihubGatewayBin string
	consulRunner     *consulrunner.ClusterRunner
)

var _ = BeforeSuite(func() {
	ApihubAPIBin, ApihubGatewayBin = test_helpers.BuildArtifacts()

	consulRunner = consulrunner.NewClusterRunner(
		9201+GinkgoParallelNode()*consulrunner.PortOffsetLength,
		1,
		"http",
	)
	consulRunner.Start()
	consulRunner.WaitUntilReady()
})

var _ = AfterSuite(func() {
	test_helpers.CleanArtifacts()
	consulRunner.Stop()
})

var _ = BeforeEach(func() {
	Expect(consulRunner.Reset()).To(Succeed())
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(time.Second * 10)
	RunSpecs(t, "Integration Suite")
}
