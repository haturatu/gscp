package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gscp/internal/config"
	"gscp/internal/local"
	"gscp/internal/remote"
	"gscp/internal/transfer"
)

func main() {
	// Parse command-line arguments
	cfg := config.ParseArgs()

	if cfg.Verbose {
		log.Printf("Starting parallel SCP with configuration:")
		log.Printf("  Destination: %s", cfg.Dest)
		log.Printf("  Remote User: %s", cfg.RemoteUser)
		log.Printf("  Remote Host: %s", cfg.RemoteHost)
		log.Printf("  Remote Path: %s", cfg.RemotePath)
		log.Printf("  Parallelism: %d", cfg.Parallelism)
		log.Printf("  Lock Dir: %s", cfg.LockDir)
		log.Printf("  Cipher: %s", cfg.CipherOption)
	}

	// Create lock directory if it doesn't exist
	err := os.MkdirAll(cfg.LockDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create lock directory: %v", err)
	}

	// Get remote file list
	remoteFiles, err := remote.GetFiles(cfg)
	if err != nil {
		log.Fatalf("Failed to get remote files: %v", err)
	}
	if cfg.Verbose {
		log.Printf("Found %d remote files", len(remoteFiles))
	}

	// Get local file list
	localFiles, err := local.GetFiles(cfg.Dest)
	if err != nil {
		log.Fatalf("Failed to get local files: %v", err)
	}
	if cfg.Verbose {
		log.Printf("Found %d local files", len(localFiles))
	}

	// Compute diff (files that exist on remote but not on local)
	diffFiles := transfer.GetDiffFiles(remoteFiles, localFiles)
	if cfg.Verbose {
		log.Printf("Found %d files to copy", len(diffFiles))
	}

	// If only listing files, print them and exit
	if cfg.OnlyListFiles {
		for _, file := range diffFiles {
			fmt.Println(file)
		}
		return
	}

	// Copy files in parallel
	start := time.Now()
	copiedFiles := transfer.CopyFilesInParallel(diffFiles, cfg)
	elapsed := time.Since(start)
	log.Printf("Done! Copied %d/%d files in %s", copiedFiles, len(diffFiles), elapsed)
}
