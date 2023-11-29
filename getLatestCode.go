package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const title = "Team Shrec code updater"

func main() {
	greeting(title)
	targetDirName := "SRC"
	targetDirValue, exists := os.LookupEnv(targetDirName)
	if !exists {
		fmt.Printf("environment variable %s is not set\n", targetDirName)
	}
	targetDir := targetDirValue
	fmt.Printf("The %v directory will be parsed for git repos and the latest code pulled for each.\n\n", targetDir)
	pullLatestCode(targetDir)
}

func pullLatestCode(targetDir string) {
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		// //skip .terraform directories
		if info.IsDir() && info.Name() == ".terraform" {
			return filepath.SkipDir
		}

		// Check if the directory is a git repo
		if info.IsDir() && isGitRepository(path) {
			fmt.Printf("Processing Git repo: %s\n", path)

			// Check if a remote has been configured
			remoteURL, err := getGitRemote(path)
			if err != nil {
				fmt.Printf("Error checking for remote for %s: %v\n", path, err)
				return nil
			}
			if remoteURL == "" {
				fmt.Printf("Skipping %v, as no remote has been defined.\n\n", path)
				return nil
			}

			//Pull the latest code
			err = gitPull(path)
			if err != nil {
				fmt.Printf("error pulling the latest code for %s: %v\n", path, err)
				return nil
			}

		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", targetDir, err)
	}
}

func isGitRepository(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func getGitRemote(path string) (string, error) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	remoteLines := strings.Split(string(output), "\n")
	if len(remoteLines) == 0 {
		return "", nil
	}
	remoteParts := strings.Fields(remoteLines[0])
	if len(remoteParts) == 0 {
		return "", nil
	}
	return remoteParts[1], nil
}

func gitPull(path string) error {
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = path
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return err
	}

	if len(statusOutput) > 0 {
		fmt.Printf("Skipping pull for %s: There are uncommitted Changes:\n%s\n", path, statusOutput)
		return nil
	}

	pullCmd := exec.Command("git", "pull")
	pullCmd.Dir = path
	err = pullCmd.Run()
	return err
}
