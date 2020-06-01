package exporter

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

// ParseEnvVar provides a wrapper to attempt to fetch the specified env var
func ParseEnvVar(envVarName string, decode bool) (string, error) {
	value := os.Getenv(envVarName)
	if decode {
		v, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return "", errors.Errorf("error decoding environment variable %q", envVarName)
		}
		value = fmt.Sprintf("%s", v)
	}
	return value, nil
}

// ExecuteCommand executes a command
func ExecuteCommand(logErr bool, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if logErr {
			klog.Errorf("%s failed output is:\n", command)
			klog.Errorf("%s\n", string(output))
		}
		return output, errors.Wrapf(err, "%s execution failed", command)
	}
	return output, nil
}
