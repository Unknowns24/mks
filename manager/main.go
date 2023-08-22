package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/unknowns24/mks/utils"
)

func GenerateMicroservice(serviceName string, features []string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	basePath := filepath.Join(currentDir, serviceName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	createBaseFiles(basePath)

	if utils.SliceContainsElement(features, "all") {
		err := AddAllFeatures(basePath)
		if err != nil {
			return err
		}
	} else {
		for _, feature := range features {
			err := AddFeature(basePath, feature)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Microservice '%s' with features %v generated successfully!\n", serviceName, features)
	return nil
}

func createBaseFiles(basePath string) {
	// Implement creation of basic files (main.go, Dockerfile, etc.)
}
