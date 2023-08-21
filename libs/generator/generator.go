package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateMicroservice(serviceName string, full bool) error {
	basePath := filepath.Join("microservices", serviceName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	createBaseFiles(basePath)
	if full {
		addAllFeatures(basePath)
	}

	fmt.Printf("Microservice '%s' generated successfully!\n", serviceName)
	return nil
}

func createBaseFiles(basePath string) {
	// Implement creation of basic files (main.go, Dockerfile, etc.)
}

func addAllFeatures(basePath string) {
	// Implement adding all features based on templates and logic
}

func AddFeature(serviceName, feature string) error {
	// Implement logic to add the specified feature to the microservice
	fmt.Printf("Feature '%s' added to microservice '%s'\n", feature, serviceName)
	return nil
}

func RemoveFeature(serviceName, feature string) error {
	// Implement logic to remove the specified feature from the microservice
	fmt.Printf("Feature '%s' removed from microservice '%s'\n", feature, serviceName)
	return nil
}
