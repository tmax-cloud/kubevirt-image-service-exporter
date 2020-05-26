package exporter

import (
	"io/ioutil"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

// MockDestination
type MockDestination struct {
	transferResponse ProcessingPhase
}

// Transfer is called to transfer the image to the Destination
func (md *MockDestination) Transfer() (ProcessingPhase, error) {
	if md.transferResponse == ProcessingPhaseError {
		return ProcessingPhaseError, errors.New("Transfer errored")
	}
	return md.transferResponse, nil
}

// Close closes any senders or other open resources.
func (md *MockDestination) Close() {
}

var _ = Describe("NewProcessor", func() {
	It("should return correct initial Processor", func() {
		md := &MockDestination{
			transferResponse: ProcessingPhaseComplete,
		}
		p := NewProcessor(md, "source/disk.img", "export", "export/disk.img")
		Expect(p.currentPhase).To(Equal(ProcessingPhaseConvert))
		Expect(p.destination).To(Equal(md))
		Expect(p.sourcePath).To(Equal("source/disk.img"))
		Expect(p.exportDir).To(Equal("export"))
		Expect(p.exportPath).To(Equal("export/disk.img"))
	})
})

var _ = Describe("Process", func() {
	BeforeEach(func() {
		err := os.MkdirAll("export", os.ModePerm)
		Expect(err).ToNot(HaveOccurred())
	})
	table.DescribeTable("Processing", func(execFunc execFunctionType, dest string, phase ProcessingPhase, errString string) {
		replaceExecFunction(execFunc, func() {
			md := &MockDestination{
				transferResponse: phase,
			}
			p := NewProcessor(md, "source/disk.img", dest, dest+"/disk.img")
			err := p.Process()
			if errString == "" {
				Expect(err).NotTo(HaveOccurred())
				Expect(p.currentPhase).To(Equal(ProcessingPhaseComplete))
			} else {
				Expect(err).To(HaveOccurred())
				Expect(p.currentPhase).To(Equal(phase))
				Expect(strings.Contains(err.Error(), errString)).To(BeTrue())
			}
		})
	},
		table.Entry("should return error when clean directory fail",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "wrong", ProcessingPhaseConvert, "Failure cleaning up export space"),
		table.Entry("should return error on convert phase",
			mockExecFunction("exit 1", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "export", ProcessingPhaseError, "Unable to convert source data to Qcow2 format"),
		table.Entry("should return error on transfer phase",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "export", ProcessingPhaseError, "Unable to transfer Qcow2 image to destination"),
		table.Entry("should return error on unknown phase",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "export", ProcessingPhase("unknown"), "Unknown processing phase"),
		table.Entry("should success and phase is ProcessingPhaseComplete",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "export", ProcessingPhaseComplete, ""),
	)
	AfterEach(func() {
		err := os.RemoveAll("export")
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = Describe("convert", func() {
	table.DescribeTable("Converting", func(execFunc execFunctionType, errString string) {
		replaceExecFunction(execFunc, func() {
			md := &MockDestination{
				transferResponse: ProcessingPhaseComplete,
			}
			p := NewProcessor(md, "source/disk.img", "export", "export/disk.img")
			phase, err := p.convert()
			if errString == "" {
				Expect(err).NotTo(HaveOccurred())
				Expect(phase).To(Equal(ProcessingPhaseTransfer))
			} else {
				Expect(err).To(HaveOccurred())
				Expect(phase).To(Equal(ProcessingPhaseError))
				Expect(strings.Contains(err.Error(), errString)).To(BeTrue())
			}
		})
	},
		table.Entry("should success and return ProcessingPhaseTransfer",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), ""),
		table.Entry("should fail and return error and ProcessingPhaseError when convert fail",
			mockExecFunction("exit 1", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "Conversion to Qcow2 failed"),
	)
})

var _ = Describe("cleanDir", func() {
	BeforeEach(func() {
		err := os.MkdirAll("source", os.ModePerm)
		Expect(err).ToNot(HaveOccurred())
		_, err = os.Create("source/disk.img")
		Expect(err).ToNot(HaveOccurred())
	})

	table.DescribeTable("Clean Directory contents", func(dest string, success bool) {
		err := cleanDir(dest)
		if success {
			Expect(err).ToNot(HaveOccurred())
			dir, err2 := ioutil.ReadDir(dest)
			Expect(err2).NotTo(HaveOccurred())
			Expect(0).To(Equal(len(dir)))
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		table.Entry("should clean directory when directory is accessible", "source", true),
		table.Entry("should return error in non existing directory", "wrong", false),
	)
	AfterEach(func() {
		err := os.RemoveAll("source")
		Expect(err).ToNot(HaveOccurred())
	})
})
