package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Schedule   string            `yaml:"schedule"`   // Cron schedule expression
	Sources    []string          `yaml:"sources"`    // Source directories to backup
	Destination Destination      `yaml:"destination"`
	Retention  RetentionPolicy   `yaml:"retention"`
	Compression CompressionType  `yaml:"compression,omitempty"`
}

// Destination represents where backups should be stored
type Destination struct {
	Type   DestinationType `yaml:"type"`   // "local" or "s3"
	Path   string          `yaml:"path"`   // Local path or S3 bucket name
	Prefix string          `yaml:"prefix,omitempty"` // Optional prefix for backup files

	// S3-specific configuration
	S3Config *S3Config `yaml:"s3,omitempty"`
}

// S3Config contains S3-specific settings
type S3Config struct {
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint,omitempty"`
	AccessKeyID     string `yaml:"access_key_id,omitempty"`
	SecretAccessKey string `yaml:"secret_access_key,omitempty"`
}

// DestinationType represents the type of destination
type DestinationType string

const (
	DestinationLocal DestinationType = "local"
	DestinationS3    DestinationType = "s3"
)

// CompressionType represents the compression algorithm
type CompressionType string

const (
	CompressionGzip  CompressionType = "gzip"
	CompressionZstd  CompressionType = "zstd"
	CompressionNone  CompressionType = "none"
)

// RetentionPolicy defines the Grandfather-Father-Son retention strategy
type RetentionPolicy struct {
	Daily   int `yaml:"daily"`   // Number of daily backups to keep
	Weekly  int `yaml:"weekly"`  // Number of weekly backups to keep
	Monthly int `yaml:"monthly"` // Number of monthly backups to keep
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Schedule == "" {
		return fmt.Errorf("schedule is required")
	}

	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source directory is required")
	}

	if c.Destination.Type == "" {
		return fmt.Errorf("destination type is required")
	}

	if c.Destination.Path == "" {
		return fmt.Errorf("destination path is required")
	}

	if c.Destination.Type != DestinationLocal && c.Destination.Type != DestinationS3 {
		return fmt.Errorf("destination type must be 'local' or 's3'")
	}

	if c.Destination.Type == DestinationS3 && c.Destination.S3Config == nil {
		return fmt.Errorf("s3 configuration is required when destination type is 's3'")
	}

	if c.Retention.Daily < 0 || c.Retention.Weekly < 0 || c.Retention.Monthly < 0 {
		return fmt.Errorf("retention values must be non-negative")
	}

	return nil
}
