package manager

import (
	"fmt"
	"os"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/utils"
)

func IsValidFeature(feature string) bool {
	valid := false

	for _, validFeature := range config.Features {
		if feature == validFeature {
			valid = true
			break
		}
	}

	return valid
}

func AddFeature(feature string) error {
	var err error

	// If global variable basePath is empty fill it
	if basePath == "" {
		// Get the current working directory
		basePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// If global variable serviceName is empty fill it
	if serviceName == "" {
		// Get Mircoservice module name
		serviceName, err = utils.GetThisModuleName()
		if err != nil {
			return err
		}
	}

	switch feature {
	case config.FEATURE_JWT:
		return addJWTFeature()
	case config.FEATURE_MYSQL:
		return addMySQLFeature()
	case config.FEATURE_GRPC_SERVER:
		return addGrpcServerFeature()
	case config.FEATURE_GRPC_CLIENT:
		return addGrpcClientFeature()
	case config.FEATURE_RMQ_PRODUCER:
		return addRMQProducerFeature()
	case config.FEATURE_RMQ_CONSUMER:
		return addRMQConsumerFeature()
	default:
		return fmt.Errorf("unknown feature: %s", feature)
	}
}

func AddAllFeatures() error {
	// Implement logic to add MySQL feature
	// Prompt the user for host, user, password, and database details
	return nil
}

func addRMQProducerFeature() error {
	// Implement logic to add RMQ feature
	// You can use templates or generate necessary files here
	return nil
}

func addRMQConsumerFeature() error {
	// Implement logic to add RMQ feature
	// You can use templates or generate necessary files here
	return nil
}

func addGrpcServerFeature() error {
	// Implement logic to add gRPC feature
	// You can use templates or generate necessary files here
	return nil
}

func addGrpcClientFeature() error {
	// Implement logic to add gRPC feature
	// You can use templates or generate necessary files here
	return nil
}

func addMySQLFeature() error {
	// Implement logic to add MySQL feature
	// Prompt the user for host, user, password, and database details
	return nil
}

func addJWTFeature() error {
	// Implement logic to add JWT feature
	// Prompt the user for JWT configuration details
	return nil
}
