# MKS CLI

mks is a command-line interface (CLI) tool that helps you generate basic microservices structures with different features. It allows you to quickly set up microservices projects by creating the necessary files and directories based on your requirements.

## Installation

To install mks, you need to have [Go](https://go.dev/doc/install) installed on your system. Then, run the following command:

```sh
go install github.com/unknowns24/mks
```

## Usage

### Build a Microservice

To create a basic microservice structure, use the `build` command:

```sh
mks build [name] --features=[features]

```

-   `name`: The name of the microservice to be generated.
-   `features`: (Optional) Generate microservice with features.

_Example:_

```sh
mks build ms_apps --features="mysql,jwt"

```

### Add a Feature

To add a specific feature to an existing microservice, use the `add` command:

```sh
mks add [feature]

```

-   `feature`: The feature you want to add.

_Example:_

```sh
mks add mysql

```

## Features

mks supports the following features:

-   MySQL Database (mysql)
-   gRPC Client (grpc-cl)
-   gRPC Server (grpc-sv)
-   RabbitMQ Producer (rmq-producer)
-   RabbitMQ Consumer (rmq-consumer)
-   Jason Web Token Auth (jwt)
-   All features (all)

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! If you find any issues or would like to add new features, feel free to open an issue or submit a pull request on the GitHub repository: https://github.com/unknowns24/mks.

## Acknowledgments

-   The CLI is built using the Cobra library: https://github.com/spf13/cobra
