package transfer

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"gscp/internal/config"
)

// GetDiffFiles computes which files exist on remote but not on local
func GetDiffFiles(remoteFiles []string, localFiles map[string]bool) []string {
	var diffFiles []string

	for _, remoteFile := range remoteFiles {
		if !localFiles[remoteFile] {
			diffFiles = append(diffFiles, remoteFile)
		}
	}

	return diffFiles
}

// CopyFilesInParallel copies multiple files in parallel from remote to local
func CopyFilesInParallel(files []string, config config.Configuration) int {
	if len(files) == 0 {
		log.Println("No files to copy")
		return 0
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Parallelism)

	totalFiles := len(files)
	copiedFiles := 0
	var mu sync.Mutex

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := CopyFile(file, config); err != nil {
				log.Printf("Error copying file %s: %v", file, err)
			} else {
				mu.Lock()
				copiedFiles++
				if config.Verbose && copiedFiles%10 == 0 {
					log.Printf("Progress: %d/%d files copied (%.2f%%)",
						copiedFiles, totalFiles, float64(copiedFiles)/float64(totalFiles)*100)
				}
				mu.Unlock()
			}
		}(file)
	}

	wg.Wait()
	return copiedFiles
}

// CopyFile copies a single file from remote to local
func CopyFile(file string, config config.Configuration) error {
	// Create lock file path by replacing slashes with underscores
	lockFileName := strings.ReplaceAll(file, "/", "_") + ".lock"
	lockFilePath := filepath.Join(config.LockDir, lockFileName)

	// Try to create and lock the file
	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot create lock file: %v", err)
	}
	defer lockFile.Close()

	// Create destination directory
	destDir := filepath.Join(config.Dest, filepath.Dir(file))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}

	// Check if file already exists in destination
	destFile := filepath.Join(config.Dest, file)
	if _, err := os.Stat(destFile); err == nil {
		// File already exists, skip
		return nil
	}

	// Construct scp command
	src := fmt.Sprintf("%s@%s:%s/%s", config.RemoteUser, config.RemoteHost, config.RemotePath, file)
	args := []string{
		"-q", // Quiet mode
		"-c", config.CipherOption,
		src,
		destFile,
	}

	if config.Verbose {
		log.Printf("Copying file: %s", file)
	}

	cmd := exec.Command("scp", args...)

	// Capture stdout and stderr
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	// Execute command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("scp failed: %v, stderr: %s", err, stderrBuf.String())
	}

	// Clean up lock file on success
	os.Remove(lockFilePath)

	return nil
}
