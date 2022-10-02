package s3

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client returns an S3 client
func Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
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
	result, err := client.GetObject(context.TODO(), input)
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(result.Body)
}

// PutObject writes an object to a bucket
func PutObject(bucket, key, body string) error {
	client, err := Client()
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), input)
	return err
}
