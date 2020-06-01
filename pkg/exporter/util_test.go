package exporter

import (
	"encoding/base64"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseEnvVar", func() {

	BeforeEach(func() {
		err := os.Setenv("value1", "value1")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv("value2", base64.StdEncoding.EncodeToString([]byte("value2")))
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv("value3", "invalid --- *** &&&")
		Expect(err).NotTo(HaveOccurred())
	})
	table.DescribeTable("Parse Env", func(envVarName, value string, decode, success bool) {
		result, err := ParseEnvVar(envVarName, decode)
		if success {
			Expect(err).NotTo(HaveOccurred())
			if result != value {
				Fail(fmt.Sprintf("got wrong value: %s, expected %s", result, value))
			}
		} else {
			Expect(err).To(HaveOccurred())
		}

	},
		table.Entry("should parse unencoded env", "value1", "value1", false, true),
		table.Entry("should parse encoded env", "value2", "value2", true, true),
		table.Entry("should parse empty string when no env", "value4", "", false, true),
		table.Entry("should return error when parsing invalid encoded env", "value3", "", true, false),
	)
	AfterEach(func() {
		err := os.Unsetenv("value1")
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv("value2")
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv("value3")
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("ExecuteCommand", func() {
	table.DescribeTable("Execute command", func(success, logErr bool, command string, args ...string) {
		_, err := ExecuteCommand(logErr, command, args...)
		if success {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		table.Entry("should success", true, false, "sleep", "1"),
		table.Entry("should exit bad", false, false, "sleep", "-1"),
	)
})
