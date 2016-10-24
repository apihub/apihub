package api_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apihub Server Suite")
}

func stringify(data []byte) string {
	return strings.Trim(string(data), "\n")
}
