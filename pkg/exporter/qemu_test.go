package exporter

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

type execFunctionType func(bool, string, ...string) ([]byte, error)

var _ = Describe("Convert to Qcow", func() {
	BeforeEach(func() {
		_, err := os.Create("dest")
		Expect(err).NotTo(HaveOccurred())
	})

	table.DescribeTable("Converting", func(execFunc execFunctionType, source, dest, errString string) {
		replaceExecFunction(execFunc, func() {
			err := convertToQcow2(source, dest)
			if errString == "" {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errString)).To(BeTrue())
			}
		})
	},
		table.Entry("should return success if qemu-img exec function returns no error",
			mockExecFunction("", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "source", "dest", ""),
		table.Entry("should return error if qemu-img exec function returns error",
			mockExecFunction("exit 1", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "source", "dest", "could not convert image to qcow2"),
		table.Entry("should return error and fail to delete dest file if qemu-img exec function returns error and dest file can't delete",
			mockExecFunction("exit 1", "convert", "-t", "none", "-p", "-c", "-O", "qcow2"), "source", "aaaa", "fail to remove aborted image"),
	)
	AfterEach(func() {
		err := os.RemoveAll("dest")
		Expect(err).NotTo(HaveOccurred())
	})
})

func mockExecFunction(errString string, checkArgs ...string) execFunctionType {
	return func(logErr bool, cmd string, args ...string) (bytes []byte, err error) {
		for _, ca := range checkArgs {
			found := false
			for _, a := range args {
				if ca == a {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		}

		bytes = []byte("")
		if errString != "" {
			err = errors.New(errString)
		}

		return
	}
}

func replaceExecFunction(replacement execFunctionType, f func()) {
	orig := qemuExecFunction
	if replacement != nil {
		qemuExecFunction = replacement
		defer func() { qemuExecFunction = orig }()
	}
	f()
}
