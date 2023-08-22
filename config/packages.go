package config

var PACKAGE_BASE = [...]string{"github.com/gofiber/fiber/v2", "github.com/spf13/viper"}
var PACKAGE_MYSQL = [...]string{"gorm.io/gorm", "gorm.io/driver/mysql"}
var PACKAGE_RMQ_PRODUCER = [...]string{"github.com/cenkalti/backoff/v4", "github.com/rabbitmq/amqp091-go"}

const (
	PACKAGE_JWT          = "github.com/golang-jwt/jwt/v4"
	PACKAGE_GRPC         = "google.golang.org/grpc"
	PACKAGE_RMQ_CONSUMER = "github.com/rabbitmq/amqp091-go"
)
