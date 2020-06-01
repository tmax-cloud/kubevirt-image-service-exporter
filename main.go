package main

import (
	"flag"
	"kubevirt-image-service-exporter/pkg/exporter"
	"os"

	"k8s.io/klog"
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
	// ExporterDestinationLocal indicates Destaination to export is local
	ExporterDestinationLocal = "local"
	// ExporterDestinationS3 indicates Destaination to export is s3
	// ExporterDestinationS3 = "s3"
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
	// if S3, need additional env variable

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
	// todo S3
	default:
		klog.Errorf("Unknown destination type %s\n", dest)
		if err != nil {
			klog.Errorf("%+v", err)
		}
		exitCode = 1
		return
	}
	defer destType.Close()
	processor := exporter.NewProcessor(destType, sourcePath, exportDir, exportPath)
	err = processor.Process()
	if err != nil {
		klog.Errorf("%+v", err)
		exitCode = 1
		return
	}
	klog.V(1).Infoln("Export complete")
}
