package utils

import (
	"fmt"
	"os/exec"
)

func InitGoModules(serviceName, basePath string) error {
	initCmd := exec.Command("go", "mod", "init", serviceName)
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("¨error running 'go mod init': %s\nOutput: %s", err, initOutput)
	}

	return nil
}

func InstallNeededPackages(basePath string) error {
	initCmd := exec.Command("go", "mod", "tidy")
	initCmd.Dir = basePath
	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running 'go mod tidy': %s\nOutput: %s", err, initOutput)
	}

	return nil
}
