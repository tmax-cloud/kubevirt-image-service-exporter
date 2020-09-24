package main

import (
	"flag"
	"k8s.io/klog"
	"kubevirt-image-service-exporter/pkg/exporter"
	"os"
)

const (
	// ExporterDiskImageName provides a constant for disk image name
	ExporterDiskImageName = "disk.img"
	// ExporterSourcePath  provides a constant to capture our env variable "EXPORTER_SOURCE_PATH"
	ExporterSourcePath = "EXPORTER_SOURCE_PATH"
	// ExporterExportDir  provides a constant to capture our env variable "EXPORTER_EXPORT_DIR"
	ExporterExportDir = "EXPORTER_EXPORT_DIR"
	// ExporterDestination provides a constant to capture our env variable "EXPORTER_DESTINATION"
	ExporterDestination = "EXPORTER_DESTINATION"
	// Endpoint is an endpoint of the external object storage where to export volume
	Endpoint = "ENDPOINT"
	// AccessKeyID is one of AWS-style credential which is needed when export volume to external object storage
	AccessKeyID = "AWS_ACCESS_KEY_ID"
	// SecretAccessKey is one of AWS-style credential which is needed when export volume to external object storage
	SecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	// ExporterDestinationLocal indicates Destination to export is local
	ExporterDestinationLocal = "local"
	// ExporterDestinationS3 indicates Destination to export is s3
	ExporterDestinationS3 = "s3"
)

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()
	defer klog.Flush()

	var err error

	klog.V(1).Infoln("Starting exporter")

	dest, _ := exporter.ParseEnvVar(ExporterDestination, false)
	sourcePath, _ := exporter.ParseEnvVar(ExporterSourcePath, false)
	exportDir, _ := exporter.ParseEnvVar(ExporterExportDir, false)

	if _, err = os.Stat(sourcePath); os.IsNotExist(err) {
		klog.Errorf("Source Image Path doesn't exist %s\n", sourcePath)
		exitCode = 1
		return
	}
	if _, err = os.Stat(exportDir); os.IsNotExist(err) {
		klog.Errorf("Export Directory doesn't exist %s\n", exportDir)
		exitCode = 1
		return
	}
	exportPath := exportDir + "/" + ExporterDiskImageName

	var destType exporter.DataDestinationInterface
	switch dest {
	case ExporterDestinationLocal:
		destType, err = exporter.NewLocalDestination()
		if err != nil {
			klog.Errorf("%+v", err)
			exitCode = 1
			return
		}
	case ExporterDestinationS3:
		keyID, _ := exporter.ParseEnvVar(AccessKeyID, false)
		accessKey, _ := exporter.ParseEnvVar(SecretAccessKey, false)
		endpoint, _ := exporter.ParseEnvVar(Endpoint, false)

		if keyID == ""{
			klog.Error("AccessKeyID is missing or blank\n")
			exitCode = 1
			return
		}
		if accessKey == ""{
			klog.Error("SecretAccessKey is missing or blank\n")
			exitCode = 1
			return
		}
		if endpoint == "" {
			klog.Error("Endpoint is missing or blank\n")
			exitCode = 1
			return
		}
		destType, err = exporter.NewS3Destination(endpoint, keyID, accessKey, exportPath)
		if err != nil {
			klog.Errorf("%+v", err)
			exitCode = 1
			return
		}
	default:
		klog.Errorf("Unknown destination type %s\n", dest)
		if err != nil {
			klog.Errorf("%+v", err)
		}
		exitCode = 1
		return
	}
	processor := exporter.NewProcessor(destType, sourcePath, exportDir, exportPath)
	err = processor.Process()
	if err != nil {
		klog.Errorf("%+v", err)
		exitCode = 1
		return
	}
	klog.V(1).Infoln("Export complete")
}
