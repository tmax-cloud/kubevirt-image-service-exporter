package exporter

import (
	"os"

	"github.com/pkg/errors"
)

var (
	qemuExecFunction = ExecuteCommand
)

// convertToQcow2 converts raw image to qcow2 image
func convertToQcow2(src, dest string) error {
	_, err := qemuExecFunction(true, "qemu-img", "convert", "-t", "none", "-p", "-c", "-O", "qcow2", src, dest)
	if err != nil {
		if err2 := os.Remove(dest); err2 != nil {
			err = errors.Wrap(err, "fail to remove aborted image")
		}
		return errors.Wrap(err, "could not convert image to qcow2")
	}
	return nil
}
