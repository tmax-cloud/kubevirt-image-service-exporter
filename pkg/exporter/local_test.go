package exporter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewLocalDestination", func() {
	It("should return LocalDestination", func() {
		destType, err := NewLocalDestination()
		Expect(err).NotTo(HaveOccurred())
		Expect(destType).To(Equal(&LocalDestination{}))
	})
})

var _ = Describe("Transfer", func() {
	It("should return complete phase", func() {
		destType, err := NewLocalDestination()
		Expect(err).NotTo(HaveOccurred())
		phase, err := destType.Transfer()
		Expect(err).NotTo(HaveOccurred())
		Expect(phase).To(Equal(ProcessingPhaseComplete))
	})
})
