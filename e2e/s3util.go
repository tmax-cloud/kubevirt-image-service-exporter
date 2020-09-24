package e2e

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"kubevirt-image-service-exporter/pkg/exporter"
	"net"
	"os"
)

const (
	// Endpoint is an endpoint of the external object storage where to export volume
	Endpoint = "ENDPOINT"
	// AccessKeyID is one of AWS-style credential which is needed when export volume to external object storage
	AccessKeyID = "AWS_ACCESS_KEY_ID"
	// SecretAccessKey is one of AWS-style credential which is needed when export volume to external object storage
	SecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	// MinioPort is a port number of minio docker container
	MinioPort = "9000"
	// BucketName is
	BucketName = "test"
	// TestKeyID is
	TestKeyID = "minio"
	// TestKey is
	TestKey = "minio123"
	// TempDir indicates
	TempDir = "tmp"
	// ContainerName indicates
	ContainerName = "kise_e2e_test_minio"
)

var hostIP = "localhost"

func getHostIP() error {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return err
	}
	defer func() {
		if cerr:=conn.Close(); cerr!=nil && err == nil {
			err = cerr
		}
	}()
	hostIP = conn.LocalAddr().(*net.UDPAddr).IP.String()
	return nil
}

func startMinio() error {
	port := MinioPort + ":" + MinioPort
	domain := "MINIO_DOMAIN=" + hostIP
	keyID := "MINIO_ACCESS_KEY=" + TestKeyID
	key := "MINIO_SECRET_KEY=" + TestKey

	minioArgsList := []string{"run", "-p", port, "-d", "-e", domain, "-e", keyID, "-e", key, "--name",
		"kise_e2e_test_minio", "minio/minio", "server", "/data"}
	_, err := exporter.ExecuteCommand(true, "docker", minioArgsList...)
	return err
}

func createBucket(minioClient *minio.Client) error {
	return minioClient.MakeBucket(context.Background(), BucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false})
}

func getObject(minioClient *minio.Client) error {
	object, err := minioClient.GetObject(context.Background(), BucketName, ExporterDiskImageName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	localFile, err := os.Create("/"+TempDir+"/"+ExporterDiskImageName)
	if err != nil {
		return err
	}
	if _, err = io.Copy(localFile, object); err != nil {
		return err
	}
	return nil
}

func createS3Client() (*minio.Client, error) {
	return minio.New(hostIP + ":" + MinioPort, &minio.Options{
		Creds:  credentials.NewStaticV4(TestKeyID, TestKey, ""),
		Secure: false,
	})
}
