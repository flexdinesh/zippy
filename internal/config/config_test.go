package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid local config",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Type: DestinationLocal,
					Path: "/backups",
				},
				Retention: RetentionPolicy{
					Daily:   7,
					Weekly:  4,
					Monthly: 12,
				},
			},
			wantErr: false,
		},
		{
			name: "valid s3 config",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Type: DestinationS3,
					Path: "my-bucket",
					S3Config: &S3Config{
						Region: "us-east-1",
					},
				},
				Retention: RetentionPolicy{
					Daily:   7,
					Weekly:  4,
					Monthly: 12,
				},
			},
			wantErr: false,
		},
		{
			name: "missing schedule",
			config: Config{
				Sources: []string{"/data"},
				Destination: Destination{
					Type: DestinationLocal,
					Path: "/backups",
				},
			},
			wantErr: true,
		},
		{
			name: "missing sources",
			config: Config{
				Schedule: "0 2 * * *",
				Destination: Destination{
					Type: DestinationLocal,
					Path: "/backups",
				},
			},
			wantErr: true,
		},
		{
			name: "missing destination type",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Path: "/backups",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid destination type",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Type: "invalid",
					Path: "/backups",
				},
			},
			wantErr: true,
		},
		{
			name: "s3 without config",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Type: DestinationS3,
					Path: "my-bucket",
				},
			},
			wantErr: true,
		},
		{
			name: "negative retention",
			config: Config{
				Schedule: "0 2 * * *",
				Sources:  []string{"/data"},
				Destination: Destination{
					Type: DestinationLocal,
					Path: "/backups",
				},
				Retention: RetentionPolicy{
					Daily: -1,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	validConfig := `schedule: "0 2 * * *"
sources:
  - /data
  - /app
destination:
  type: local
  path: /backups
  prefix: backup-
retention:
  daily: 7
  weekly: 4
  monthly: 12
compression: gzip
`

	if err := os.WriteFile(configPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Schedule != "0 2 * * *" {
		t.Errorf("Expected schedule '0 2 * * *', got '%s'", cfg.Schedule)
	}

	if len(cfg.Sources) != 2 {
		t.Errorf("Expected 2 sources, got %d", len(cfg.Sources))
	}

	if cfg.Destination.Type != DestinationLocal {
		t.Errorf("Expected destination type 'local', got '%s'", cfg.Destination.Type)
	}

	if cfg.Retention.Daily != 7 {
		t.Errorf("Expected daily retention 7, got %d", cfg.Retention.Daily)
	}
}

func TestLoadInvalidFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	invalidYAML := `this is not: [valid yaml`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML")
	}
}
