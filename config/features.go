package config

const ALL_FEATURES = "all"

var Features = [...]string{
	"all",
	"jwt",
	"mysql",
	"rmq-consumer",
	"rmq-producer",
	"grpc-cl",
	"grpc-sv",
}
