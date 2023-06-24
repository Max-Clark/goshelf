package goshelf

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flags", func() {
	Describe("Parsing flags", func() {
		Context("no arguments passed", func() {
			var mockNoArgs = []string{}

			It("should return defaults", func() {
				cfg, err := InitFlags(mockNoArgs)

				Expect(err).To(BeNil())
				Expect(cfg).ToNot(BeNil())

				Expect(cfg.Host).To(Equal("0.0.0.0"))
			})

			// TODO: add more tests
		})
	})

	Describe("Usage", func() {
		Context("prints usage to stderr", func() {
			var getUsage = []string{"-h"}

			It("should print a nice usage", func() {
				cfg, err := InitFlags(getUsage)

				Expect(err).To(BeNil())
				Expect(cfg).ToNot(BeNil())

				Expect(cfg.Host).To(Equal("0.0.0.0"))
			})

			// TODO: add more tests
		})
	})
})
