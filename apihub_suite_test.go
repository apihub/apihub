package apihub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApihub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apihub Suite")
}
