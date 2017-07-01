package consuladapter_test

import (
	"os"
	"path"

	"code.cloudfoundry.org/consuladapter"
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

var (
	consulRunner           *consulrunner.ClusterRunner
	consulClient           consuladapter.Client
	err                    error
	basePath               string
	consulCACert           string
	consulClientCert       string
	consulCLientKey        string
	invalidConsulClientKey string
)

var _ = BeforeSuite(func() {
	basePath = path.Join(os.Getenv("GOPATH"), "src/code.cloudfoundry.org/consuladapter/fixtures")
	consulCACert = path.Join(basePath, "consul-ca.crt")
	consulClientCert = path.Join(basePath, "consul.crt")
	consulCLientKey = path.Join(basePath, "consul.key")
	invalidConsulClientKey = path.Join(basePath, "invalid-consul.key")

	consulRunner = consulrunner.NewClusterRunner(
		consulrunner.ClusterRunnerConfig{
			StartingPort: 9901 + config.GinkgoConfig.ParallelNode*consulrunner.PortOffsetLength,
			NumNodes:     1,
			Scheme:       "https",
			CACert:       consulCACert,
			ClientCert:   consulClientCert,
			ClientKey:    consulCLientKey,
		},
	)

	consulRunner.Start()
	consulRunner.WaitUntilReady()
	consulRunner.Reset()
})

var _ = AfterSuite(func() {
	if consulRunner != nil {
		consulRunner.Stop()
	}
})
