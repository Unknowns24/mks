package config

const ALL_FEATURES = "all"

const FEATURE_JWT = "jwt"
const FEATURE_MYSQL = "mysql"
const FEATURE_GRPC_CLIENT = "grpc-client"
const FEATURE_GRPC_SERVER = "grpc-server"
const FEATURE_RMQ_PRODUCER = "rmq-producer"
const FEATURE_RMQ_CONSUMER = "rmq-consumer"

var Features = [...]string{
	ALL_FEATURES,
	FEATURE_JWT,
	FEATURE_MYSQL,
	FEATURE_GRPC_CLIENT,
	FEATURE_GRPC_SERVER,
	FEATURE_RMQ_PRODUCER,
	FEATURE_RMQ_CONSUMER,
}
