package s3

import (
	"fmt"
	"os"
	"time"

	"github.com/ghodss/yaml"
)

// ConfigFile wraps a config from an S4 object
type ConfigFile struct {
	Bucket      string
	Key         string
	Config      interface{}
	LastUpdated int64
	OnSuccess   func(*ConfigFile)
	OnError     func(*ConfigFile, error)
}

// Load downloads and parses the S3 config object
func (c *ConfigFile) Load() error {
	if c.Bucket == "" {
		return fmt.Errorf("bucket not provided")
	}
	if c.Key == "" {
		return fmt.Errorf("s3 key not provided")
	}
	obj, err := GetObject(c.Bucket, c.Key)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(obj, c.Config)
}

// GetConfig loads a config struct from an S3 object
func GetConfig(bucket, key string, config interface{}) (*ConfigFile, error) {
	c := &ConfigFile{
		Bucket: bucket,
		Key:    key,
		Config: config,
	}
	err := c.Load()
	return c, err
}

// GetConfigFromEnv loads a config struct using bucket/object names from the environment
func GetConfigFromEnv(config interface{}) (*ConfigFile, error) {
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")
	return GetConfig(bucket, key, config)
}

// Autoreload enables reloading every $interval seconds
func (c *ConfigFile) Autoreload(delay int) {
	go func(c *ConfigFile, delay int) {
		for {
			now := time.Now().Unix()
			if c.LastUpdated+int64(delay) < now {
				if err := c.Load(); err == nil {
					if c.OnSuccess != nil {
						c.OnSuccess(c)
					}
					c.LastUpdated = time.Now().Unix()
				} else if c.OnError != nil {
					c.OnError(c, err)
				}
			}
			time.Sleep(time.Second)
		}
	}(c, delay)
}
