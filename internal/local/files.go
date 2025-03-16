package local

import (
	"os"
	"path/filepath"
)

// GetFiles retrieves a map of files in the local destination directory
func GetFiles(destDir string) (map[string]bool, error) {
	files := make(map[string]bool)

	err := filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Get relative path
			relPath, err := filepath.Rel(destDir, path)
			if err != nil {
				return err
			}
			files[relPath] = true
		}

		return nil
	})

	return files, err
}
