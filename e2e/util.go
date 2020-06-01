package e2e

import (
	"encoding/json"
	"kubevirt-image-service-exporter/pkg/exporter"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

// ImgInfo contains the virtual image information.
type ImgInfo struct {
	// Format contains the format of the image
	Format string `json:"format"`
	// BackingFile is the file name of the backing file
	BackingFile string `json:"backing-filename"`
	// VirtualSize is the disk size of the image which will be read by vm
	VirtualSize int64 `json:"virtual-size"`
	// ActualSize is the size of the qcow2 image
	ActualSize int64 `json:"actual-size"`
}

func vailidateImage(imagePath, format string) error {
	var output []byte
	var err error

	output, err = exporter.ExecuteCommand(false, "qemu-img", "info", "--output=json", imagePath)
	if err != nil {
		return errors.Wrapf(err, "Error getting info on image %s", imagePath)
	}
	var info ImgInfo
	err = json.Unmarshal(output, &info)
	if err != nil {
		klog.Errorf("Invalid JSON:\n%s\n", string(output))
		return errors.Wrapf(err, "Invalid json for image %s", imagePath)
	}
	if info.Format != format {
		return errors.Wrapf(err, "Invalid Format: image %s format is %s", imagePath, info.Format)
	}
	return nil
}
