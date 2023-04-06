package tryst_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTryst(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tryst Suite")
}
