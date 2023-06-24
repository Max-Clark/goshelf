package goshelf

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoshelf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goshelf Suite")
}
