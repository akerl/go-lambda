package s3

import (
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client returns an S3 client
func Client() (*s3.Client, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return s3.New(cfg), nil
}

// GetObject returns an object from a bucket
func GetObject(bucket, key string) ([]byte, error) {
	client, err := Client()
	if err != nil {
		return []byte{}, err
	}
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	req := client.GetObjectRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(result.Body)
}
