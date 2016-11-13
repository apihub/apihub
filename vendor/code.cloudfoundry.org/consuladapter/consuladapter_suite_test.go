package consuladapter_test

import (
	"code.cloudfoundry.org/consuladapter/consulrunner"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConsulAdapter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adapter <-> Cluster-Runner Integration Suite")
}

const clusterSize = 1

var clusterRunner *consulrunner.ClusterRunner

var _ = BeforeSuite(func() {
	clusterStartingPort := 5001 + config.GinkgoConfig.ParallelNode*consulrunner.PortOffsetLength*clusterSize
	clusterRunner = consulrunner.NewClusterRunner(clusterStartingPort, clusterSize, "http")
})
