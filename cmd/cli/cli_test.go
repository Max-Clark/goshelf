package cli

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Suite")
}

// Mock structure & Read for testing
type R struct{}

func (r *R) Read(p []byte) (n int, err error) {
	message := "getting input from user\n"
	for i := 0; i < len(message); i++ {
		p[i] = message[i]
	}
	return len(message), nil
}

var _ = Describe("Cli", func() {
	Context("getting input from user", func() {
		r := &R{}

		It("should get a valid string", func() {
			prompt := "test prompt"
			GetCliPrompt(&prompt, r)
		})

		// TODO: add more tests
	})

})
