package e2e

import (
	"kubevirt-image-service-exporter/pkg/exporter"
	"os"

	. "github.com/onsi/ginkgo" //nolint // use ginkgo
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega" //nolint // use gomega
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

	// SourceDir indicates Source image directory
	SourceDir = "source"
	// ExportDir indicates Source image directory
	ExportDir = "export"
	// SourcePath indicates Source image path
	SourcePath = SourceDir + "/" + ExporterDiskImageName
	// ExportPath indicates Export image path
	ExportPath = ExportDir + "/" + ExporterDiskImageName
	// TempSourceImagePath indicates temp source image path
	TempSourceImagePath = SourceDir + "/temp.img"
	// ExportFormat indicates export image path
	ExportFormat = "qcow2"
	// Exporter indicates exporter program name
	Exporter = "kubevirt-image-service-exporter"
)

var _ = Describe("Test Exporter", func() {

	table.DescribeTable("Exporting should", func(imageURL string, success bool) {
		processExporter(imageURL, func() {
			_, err := exporter.ExecuteCommand(false, Exporter)
			By("End Exporter")
			if success {
				Expect(err).NotTo(HaveOccurred())
				By("Validating image foramt")
				err = vailidateImage(ExportPath, ExportFormat)
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
		})
	},
		table.Entry("should return success", "https://download.cirros-cloud.net/contrib/0.3.0/cirros-0.3.0-i386-disk.img", true),
	)
	AfterEach(func() {
		err := os.RemoveAll(SourcePath)
		Expect(err).NotTo(HaveOccurred())
		err = os.RemoveAll(TempSourceImagePath)
		Expect(err).NotTo(HaveOccurred())
		err = os.RemoveAll(ExportPath)
		Expect(err).NotTo(HaveOccurred())
	})
})

func processExporter(imageURL string, f func()) {
	By("Getting image from url " + imageURL)
	_, err := exporter.ExecuteCommand(false, "wget", "-O", TempSourceImagePath, imageURL)
	Expect(err).NotTo(HaveOccurred())
	By("Converting image to raw")
	_, err = exporter.ExecuteCommand(false, "qemu-img", "convert", "-t", "none", "-p", "-O", "raw", TempSourceImagePath, SourcePath)
	Expect(err).NotTo(HaveOccurred())
	By("Start Exporter")
	f()
}
