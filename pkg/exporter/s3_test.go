package exporter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"net/url"
)

// #	NewS3Destination	access-key-id	secret-access-key	endpoint    endpoint access style
// 1	X	nil
// 2	X	vaild	vaild	empty
// 3	X	vaild	vaild	none url
// 4	O	vaild	vaild	vaild	http	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name)
// 5	O	vaild	vaild	vaild	https	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name)
// 6	O	vaild	vaild	vaild	https	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name/still-key-name)
var _ = Describe("NewS3Destination", func() {
	DescribeTable("should",
		func(endpoint, accessKeyID, secretAccessKey, exportPath string, expected bool) {
			client, err := NewS3Destination(endpoint, accessKeyID, secretAccessKey, exportPath)
			if expected {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client).ShouldNot(BeNil())
				Expect(client.filePath).Should(Equal(exportPath))
			} else {
				Expect(err).Should(HaveOccurred())
				Expect(client).Should(BeNil())
			}
		},
		Entry("return error for null endpoint", nil, "Q3AM3UQ867SPQQA43P2F", "tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", false),
		Entry("return error for empty endpoint", "", "Q3AM3UQ867SPQQA43P2F", "tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", false),
		Entry("return error for invalid endpoint url", "dfdfdkfjdl-123qjdk", "Q3AM3UQ867SPQQA43P2F", "tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", false),
		Entry("return NewS3Destination", "http://bucket-name.s3.Region.amazonaws.com/key-name", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", true),
		Entry("return NewS3Destination", "https://bucket-name.s3.Region.amazonaws.com/key-name", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", true),
		Entry("return NewS3Destination", "http://bucket-name.s3.Region.amazonaws.com/key-name/still-key-name", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "/export/disk.img", true),
		)
})

// #	extract correct information		endpoint access style
// 1	O	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name)
// 2	O	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name)
// 3	O	virtual host style (https://bucket-name.s3.Region.amazonaws.com/key-name/still-key-name)
// 4	X	path style (https://s3.Region.amazonaws.com/bucket-name/key-name)
// 5	X	path style (https://s3.Region.amazonaws.com/bucket-name/key-name/still-key-name)
var _ = Describe("splitEndpoint", func() {
	DescribeTable("should",
		func(full, ep, bucket, object string, expected bool) {
			parsed, err := url.Parse(full)
			Expect(err).ShouldNot(HaveOccurred())
			endpoint, bucketName, objectName := splitEndpoint(parsed)
			if expected {
				Expect(endpoint).Should(Equal(ep))
				Expect(bucketName).Should(Equal(bucket))
				Expect(objectName).Should(Equal(object))
			} else {
				Expect(endpoint).ShouldNot(Equal(ep))
				Expect(bucketName).ShouldNot(Equal(bucket))
				Expect(objectName).ShouldNot(Equal(object))
			}
		},
		Entry("extract correct information", "http://bucket-name.s3.Region.amazonaws.com/key-name", "s3.Region.amazonaws.com", "bucket-name", "key-name", true),
		Entry("extract correct information", "https://bucket-name.s3.Region.amazonaws.com/key-name", "s3.Region.amazonaws.com", "bucket-name", "key-name", true),
		Entry("extract correct information", "https://bucket-name.s3.Region.amazonaws.com/key-name/still-key-name", "s3.Region.amazonaws.com", "bucket-name", "key-name/still-key-name", true),
		Entry("extract wrong information", "https://s3.Region.amazonaws.com/bucket-name/key-name", "s3.Region.amazonaws.com", "bucket-name", "key-name", false),
		Entry("extract wrong information", "https://s3.Region.amazonaws.com/bucket-name/key-name/still-key-name", "s3.Region.amazonaws.com", "bucket-name", "key-name/still-key-name", false),
	)
})

// #	isSSLEnabled 	endpoint
// 1	X	http
// 2	O	https
var _ = Describe("isSSLEnabled", func() {
	DescribeTable("should",
		func(scheme string, expected bool) {
			Expect(isSSLEnabled(scheme)).Should(Equal(expected))
		},
		Entry("be false", "http", false),
		Entry("be true", "https", true),
	)
})
