package e2e

import (
	"encoding/json"
	"kubevirt-image-service-exporter/pkg/exporter"
	"os"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

const(
	// ExporterSourcePath  provides a constant to capture our env variable "EXPORTER_SOURCE_PATH"
	ExporterSourcePath = "EXPORTER_SOURCE_PATH"
	// ExporterExportDir  provides a constant to capture our env variable "EXPORTER_EXPORT_DIR"
	ExporterExportDir = "EXPORTER_EXPORT_DIR"
	// ExporterDestination provides a constant to capture our env variable "EXPORTER_DESTINATION"
	ExporterDestination = "EXPORTER_DESTINATION"
	// Exporter indicates exporter program name
	Exporter = "localhost:5000/kubevirt-image-service-exporter:canary"
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

func getArgsList(destination string) ([]string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "Can't get current path")
	}
	mountSourcePath := currentPath + "/" + SourceDir + ":/" + SourceDir
	mountExportPath := currentPath + "/" + ExportDir + ":/" + ExportDir
	envExporterSourcePath := ExporterSourcePath + "=/" + SourcePath
	envExporterExportDir := ExporterExportDir + "=/" + ExportDir
	envExporterDestination := ExporterDestination + "=" + destination
	argsList := []string{"run", "-v", mountSourcePath, "-v", mountExportPath, "-e", envExporterDestination, "-e", envExporterSourcePath, "-e", envExporterExportDir}

	if destination == ExporterDestinationS3 {
		endpointURL := Endpoint + "=" + "http://" + BucketName + "." + hostIP + ":" + MinioPort + "/" + ExporterDiskImageName
		keyID := AccessKeyID + "=" + TestKeyID
		key := SecretAccessKey + "=" + TestKey
		argsList = append(argsList, "-e", keyID, "-e", key, "-e", endpointURL)
	}
	return append(argsList, Exporter), nil
}
