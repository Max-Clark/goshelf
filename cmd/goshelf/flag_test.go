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
				cfg, flagSet, err := InitFlags(mockNoArgs)

				Expect(err).To(BeNil())
				Expect(flagSet).ToNot(BeNil())
				Expect(cfg).ToNot(BeNil())

				Expect(cfg.Host).To(Equal("0.0.0.0"))
			})

			// TODO: add more tests
		})

		Context("cli argument passed", func() {
			var mockNoArgs = []string{"ignored", "BookRead"}

			It("should return cli argument", func() {
				cfg, flagSet, err := InitFlags(mockNoArgs)

				Expect(err).To(BeNil())
				Expect(flagSet).ToNot(BeNil())
				Expect(cfg).ToNot(BeNil())

				noFlagArgs := flagSet.Args()
				Expect(noFlagArgs).To(ContainElement("BookRead"))
			})

			// TODO: add more tests
		})
	})
})
