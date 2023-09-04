package manager

import (
	"fmt"
	"os"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/libs/addons"
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
	if global.BasePath == "" {
		// Get the current working directory
		global.BasePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// If global variable serviceName is empty fill it
	if global.ServiceName == "" {
		// Get Mircoservice module name
		global.ServiceName, err = utils.GetThisModuleName()
		if err != nil {
			return err
		}
	}

	switch feature {
	case config.FEATURE_JWT:
		return addons.InstallJWT()
	case config.FEATURE_MYSQL:
		return addons.InstallMySQL()
	case config.FEATURE_GRPC_SERVER:
		return addons.InstallGrpcServer()
	case config.FEATURE_GRPC_CLIENT:
		return addons.InstallGrpcClient()
	case config.FEATURE_RMQ_PRODUCER:
		return addons.InstallRMQProducer()
	case config.FEATURE_RMQ_CONSUMER:
		return addons.InstallRMQConsumer()
	default: // unrechable code
		return fmt.Errorf("unknown feature: %s", feature)
	}
}

func AddAllFeatures() error {
	// Installing JWT feature
	if jwtErr := addons.InstallJWT(); jwtErr != nil {
		return jwtErr
	}

	// Installing MySQL feature
	if mysqlErr := addons.InstallMySQL(); mysqlErr != nil {
		return mysqlErr
	}

	// Installing gRPC server feature
	if grpcSvErr := addons.InstallGrpcServer(); grpcSvErr != nil {
		return grpcSvErr
	}

	// Installing gRPC client feature
	if grpcClErr := addons.InstallGrpcClient(); grpcClErr != nil {
		return grpcClErr
	}

	// Installing RabbitMQ producer feature
	if rmqProdErr := addons.InstallRMQProducer(); rmqProdErr != nil {
		return rmqProdErr
	}

	// Installing RabbitMQ consumer feature
	if rmqConsErr := addons.InstallRMQConsumer(); rmqConsErr != nil {
		return rmqConsErr
	}
	return nil
}
