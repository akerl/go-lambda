package s3

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

// Client returns an S3 client
func Client() *s3.S3 {
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))
	return s3.New(sess)
}

// GetObject returns an object from a bucket
func GetObject(bucket, key string) ([]byte, error) {
	client := Client()
	obj, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(obj.Body)
}

// GetConfig loads a config struct from an S3 object
func GetConfig(bucket, key string, opts interface{}) error {
	obj, err := GetObject(bucket, key)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(obj, opts)
}

// GetConfigFromEnv loads a config struct using bucket/object names from the environment
func GetConfigFromEnv(opts interface{}) error {
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")
	if bucket == "" {
		return fmt.Errorf("bucket not provided")
	}
	if key == "" {
		return fmt.Errorf("s3 key not provided")
	}

	return GetConfig(bucket, key, opts)
}
