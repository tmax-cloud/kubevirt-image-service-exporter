package exporter

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// SupportPathStyleRequest indicates that only virtual hosted-style bucket is supported to export
const SupportPathStyleRequest = false
// SplitNumber indicates that expected length of array is 2 after splitting bucket name from the given endpoint
const SplitNumber = 2

// S3Destination is the struct containing the information needed to export to s3 destination
type S3Destination struct{
	s3Client *minio.Client
	bucketName string
	objectName string
	filePath string
}

// NewS3Destination creates a new instance of the s3 data transfer
func NewS3Destination(endpoint, accessKeyID, secretAccessKey, exportPath string) (*S3Destination, error) {
	fullURL, err := validateEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpoint, bucket, object := splitEndpoint(fullURL)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: isSSLEnabled(fullURL.Scheme),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Initializing s3 client is failed")
	}

	return &S3Destination{
		s3Client: client,
		filePath: exportPath,
		bucketName: bucket,
		objectName: object,
	}, nil
}

// Transfer puts qcow2 image to the S3Destination
func (sd *S3Destination) Transfer() (ProcessingPhase, error) {
	_, err := sd.s3Client.FPutObject(context.Background(), sd.bucketName, sd.objectName, sd.filePath, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return ProcessingPhaseError, errors.Wrap(err, "Upload to object storage is failed")
	}
	return ProcessingPhaseComplete, nil
}

func validateEndpoint(endpoint string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme == "" {
		return nil, errors.Errorf("URI scheme is missing or blank")
	}
	if parsed.Path == "" {
		return nil, errors.Errorf("Object name is missing or blank")
	}
	if hosts := strings.SplitN(parsed.Host, ".", SplitNumber); len(hosts) < SplitNumber {
		return nil, errors.Errorf("Bucket name is missing or blank")
	}
	return parsed, nil
}

func splitEndpoint(full *url.URL) (endpoint, bucketName, objectName string) {
	if SupportPathStyleRequest {
		names := strings.SplitN(full.Path, "/", 2)
		return full.Host, names[0], names[1]
	}
	// Support virtual hosted style request
	hosts := strings.SplitN(full.Host, ".", 2)
	return hosts[1], hosts[0], strings.Trim(full.Path, "/")
}

func isSSLEnabled(scheme string) (useSSL bool) {
	return scheme == "https"
}
