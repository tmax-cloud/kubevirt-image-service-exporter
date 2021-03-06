package e2e

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	err := os.MkdirAll(SourceDir, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
	err = os.MkdirAll(ExportDir, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(SourceDir)
	Expect(err).NotTo(HaveOccurred())
	err = os.RemoveAll(ExportDir)
	Expect(err).NotTo(HaveOccurred())
})
