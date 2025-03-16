package config

import (
	"flag"
	"fmt"
	"os"
)

// Configuration holds all the parameters for the SCP operation
type Configuration struct {
	Dest          string
	RemoteUser    string
	RemoteHost    string
	RemotePath    string
	Parallelism   int
	LockDir       string
	CipherOption  string
	Verbose       bool
	OnlyListFiles bool
}

// ParseArgs parses command-line arguments and returns a Configuration
func ParseArgs() Configuration {
	config := Configuration{
		LockDir:      "/tmp/scp_lock",
		Parallelism:  10,
		CipherOption: "aes128-gcm@openssh.com",
	}

	flag.StringVar(&config.Dest, "d", "", "Destination directory (local)")
	flag.StringVar(&config.RemoteUser, "u", "", "Remote username")
	flag.StringVar(&config.RemoteHost, "h", "", "Remote hostname or IP")
	flag.StringVar(&config.RemotePath, "r", "", "Remote path")
	flag.IntVar(&config.Parallelism, "P", 10, "Number of parallel SCP processes")
	flag.StringVar(&config.LockDir, "L", "/tmp/scp_lock", "Lock directory")
	flag.StringVar(&config.CipherOption, "c", "aes128-gcm@openssh.com", "SSH cipher option")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&config.OnlyListFiles, "l", false, "Only list files to be copied, don't copy")

	flag.Parse()

	// Check required parameters
	if config.Dest == "" || config.RemoteUser == "" || config.RemoteHost == "" || config.RemotePath == "" {
		flag.Usage()
		fmt.Println("\nRequired flags: -d, -u, -h, -r")
		os.Exit(1)
	}

	return config
}
