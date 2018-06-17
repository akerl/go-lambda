package s3

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// S3Config wraps a config from an S4 object
type S3Config struct {
	Bucket string
	Key    string
	Config interface{}
}

func (s *S3Config) Load() error {
	if s.Bucket == "" {
		return fmt.Errorf("bucket not provided")
	}
	if s.Key == "" {
		return fmt.Errorf("s3 key not provided")
	}
	obj, err := GetObject(s.Bucket, s.Key)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(obj, s.Config)
}

// GetConfig loads a config struct from an S3 object
func GetConfig(bucket, key string, config interface{}) (*S3Config, error) {
	s := &S3Config{
		Bucket: bucket,
		Key:    key,
		Config: config,
	}
	err := s.Load()
	return s, err
}

// GetConfigFromEnv loads a config struct using bucket/object names from the environment
func GetConfigFromEnv(config interface{}) (*S3Config, error) {
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")
	return GetConfig(bucket, key, config)
}

// Autoreload enables reloading every $interval seconds
func (s *S3Config) Autoreload(delay int) {
	go func(s *S3Config, delay int) {
		last = time.Now()
		for {
			// TODO: Check if I need to try to update
			if err := s.Load(); err != nil {
				last = time.Now()
			}
			time.Sleep(1)
		}
	}(s, delay)
}
