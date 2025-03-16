package remote

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"

	"gscp/internal/config"
)

// GetFiles retrieves the list of files from the remote server
func GetFiles(config config.Configuration) ([]string, error) {
	if config.Verbose {
		log.Println("Getting remote file list...")
	}

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", config.RemoteUser, config.RemoteHost),
		fmt.Sprintf("find '%s' -type f | sed 's|^%s/||'", config.RemotePath, config.RemotePath))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting ssh command: %v", err)
	}

	var files []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		file := scanner.Text()
		if file != "" {
			files = append(files, file)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ssh output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ssh command failed: %v", err)
	}

	return files, nil
}
