package s3

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

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
