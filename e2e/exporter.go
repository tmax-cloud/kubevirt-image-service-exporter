package e2e

import (
	"kubevirt-image-service-exporter/pkg/exporter"
	"os"

	. "github.com/onsi/ginkgo" //nolint
	. "github.com/onsi/gomega" //nolint
)

const (
	// ExporterDiskImageName provides a constant for disk image name
	ExporterDiskImageName = "disk.img"
	// ExporterDestinationLocal indicates Destination to export is local
	ExporterDestinationLocal = "local"
	// ExporterDestinationS3 indicates Destination to export is s3
	ExporterDestinationS3 = "s3"

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
)

var _ = Describe("Test Exporter Destination", func() {
	BeforeEach(func() {
		imageURL := "https://download.cirros-cloud.net/contrib/0.3.0/cirros-0.3.0-i386-disk.img"
		By("Getting image from url " + imageURL)
		_, err := exporter.ExecuteCommand(false, "wget", "-O", TempSourceImagePath, imageURL)
		Expect(err).NotTo(HaveOccurred())
		By("Converting image to raw")
		_, err = exporter.ExecuteCommand(false, "qemu-img", "convert", "-t", "none", "-p", "-O", "raw", TempSourceImagePath, SourcePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Should export to local", func() {
		By("Exporting")
		argsList, err := getArgsList(ExporterDestinationLocal)
		Expect(err).NotTo(HaveOccurred())
		_, err = exporter.ExecuteCommand(false, "docker", argsList...)
		Expect(err).NotTo(HaveOccurred())

		By("Validating image format")
		err = vailidateImage(ExportPath, ExportFormat)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Should export to s3", func() {
		By("Creating local object storage")
		err := getHostIP()
		Expect(err).NotTo(HaveOccurred())
		err = startMinio()
		Expect(err).NotTo(HaveOccurred())
		s3client, err := createS3Client()
		Expect(err).NotTo(HaveOccurred())
		err = createBucket(s3client)
		Expect(err).NotTo(HaveOccurred())

		By("Exporting")
		argsList, err := getArgsList(ExporterDestinationS3)
		Expect(err).NotTo(HaveOccurred())
		_, err = exporter.ExecuteCommand(false, "docker", argsList...)
		Expect(err).NotTo(HaveOccurred())

		By("Validating image format")
		err = vailidateImage(ExportPath, ExportFormat)
		Expect(err).NotTo(HaveOccurred())

		By("Validating exported object")
		err = getObject(s3client)
		Expect(err).NotTo(HaveOccurred())
		err = vailidateImage("/"+TempDir+"/"+ExporterDiskImageName, ExportFormat)
		Expect(err).NotTo(HaveOccurred())

		By("Removing local object storage")
		argsList = []string{"stop", ContainerName}
		_, err = exporter.ExecuteCommand(false, "docker", argsList...)
		Expect(err).NotTo(HaveOccurred())
		argsList = []string{"rm", ContainerName}
		_, err = exporter.ExecuteCommand(false, "docker", argsList...)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(SourcePath)
		Expect(err).NotTo(HaveOccurred())
		err = os.RemoveAll(TempSourceImagePath)
		Expect(err).NotTo(HaveOccurred())
		err = os.RemoveAll(ExportPath)
		Expect(err).NotTo(HaveOccurred())
	})
})
