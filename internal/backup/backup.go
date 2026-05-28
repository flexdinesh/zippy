package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"zippy/internal/config"
)

// Service handles backup operations
type Service struct {
	cfg *config.Config
}

// NewService creates a new backup service
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// Run performs a backup operation
func (s *Service) Run() error {
	backupName := generateBackupName(time.Now())
	backupPath := filepath.Join(s.cfg.Destination.Path, backupName)

	// Create archive
	if err := s.createArchive(backupPath); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	// Clean up old backups based on retention policy
	if err := s.cleanupOldBackups(); err != nil {
		return fmt.Errorf("failed to cleanup old backups: %w", err)
	}

	return nil
}

// createArchive creates a compressed tar archive of the source directories
func (s *Service) createArchive(outputPath string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add each source directory to the archive
	for _, source := range s.cfg.Sources {
		if err := addToArchive(tarWriter, source); err != nil {
			return fmt.Errorf("failed to add %s to archive: %w", source, err)
		}
	}

	return nil
}

// addToArchive recursively adds a directory to the tar archive
func addToArchive(tw *tar.Writer, sourcePath string) error {
	return filepath.Walk(sourcePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a regular file or directory
		if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
			return nil
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header: %w", err)
		}

		// Use relative path in archive
		relPath, err := filepath.Rel(filepath.Dir(sourcePath), file)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		header.Name = relPath

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}

		// If it's a file, write the contents
		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return fmt.Errorf("failed to write file to archive: %w", err)
			}
		}

		return nil
	})
}

// cleanupOldBackups removes backups according to retention policy
func (s *Service) cleanupOldBackups() error {
	// This is a placeholder for the retention logic
	// Will be implemented in the retention package
	return nil
}

// generateBackupName creates a timestamped backup filename
func generateBackupName(t time.Time) string {
	return fmt.Sprintf("backup-%s.tar.gz", t.Format("2006-01-02-150405"))
}
