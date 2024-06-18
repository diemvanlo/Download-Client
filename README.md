# Project Guide

This guide will help you to run the `goload` project. The project is written in Go and uses a Makefile to automate various tasks such as generating code, building the project, and running the server.

## Prerequisites

- Go (version 1.16 or later)
- Docker and Docker Compose (for running the project in a containerized environment)
- golangci-lint (for linting Go code)
- buf (for generating API code)
- wire (for generating dependency injection code)

## Releases

- `v1.0.0`: Initial release of the `goload` project. This project is written in Go and includes a Makefile to automate tasks such as generating code, building the project, and running the server. It supports multiple platforms including Linux, macOS, and Windows, both for amd64 and arm64 architectures. The project structure includes API definitions, application's entry point, internal business logic and data access code, Docker Compose files for containerized environment, and built binaries. This release includes all the functionalities as described in the `README.md`.

For more detailed information about each release, please check the [releases page](https://github.com/diemvanlo/goload/releases) on GitHub.

## Running the Project

1. Clone the repository to your local machine.

2. Navigate to the project directory.

3. Run the Makefile to generate necessary code and build the project:

```bash
make all
```

This command will execute the `generate` and `build-all` tasks defined in the Makefile. The `generate` task generates necessary code for the API and dependency injection. The `build-all` task builds the project for multiple platforms including Linux, macOS, and Windows, both for amd64 and arm64 architectures.

## Project Structure

The project is structured as follows:

- `api`: Contains the protobuf definitions for the API.
- `cmd`: Contains the application's entry point.
- `internal`: Contains the application's internal code, including business logic and data access code.
- `deployments`: Contains Docker Compose files for running the project in a containerized environment.
- `build`: Contains the built binaries of the project. This directory is created when you run the `make build-all` command.
- `configs`: Contains the configuration files for the project.

## Configuration Guide

The `goload` project uses a configuration file to manage various settings. The configuration file is named `local.yaml` and is located in the `configs` directory.

Here's a brief description of each configuration option:

- `database`: This section contains settings for the database connection. You can specify the `host`, `port`, `username`, `password`, and `database` name.

- `cache`: This section contains settings for the cache. You can specify the `type` of cache (e.g., "redis"), the `address` of the cache server, and the `username` and `password` if required.

- `mq`: This section contains settings for the message queue. You can specify the `addresses` of the message queue servers and the `client_id`.

- `auth`: This section contains settings for authentication. You can specify the `cost` for hashing passwords and the `expires_in` and `regenerate_token_before_expiry` settings for tokens.

- `grpc`: This section contains settings for the gRPC server. You can specify the `address` on which the server will listen.

- `http`: This section contains settings for the HTTP server. You can specify the `address` on which the server will listen.

- `Download`: This section contains settings for the download functionality. You can specify the `model` (e.g., "local") and the `download_dir`.

Please modify these settings as per your requirements. If you're running the project in a containerized environment using Docker Compose, you might need to adjust these settings to match your Docker Compose configuration.

## Additional Commands

- To clean the build directory, run:

```bash
make clean
```

- To run the server, use:

```bash
make run-server
```

- To bring up the Docker Compose development environment, use:

```bash
make docker-compose-dev-up
```

- To bring down the Docker Compose development environment, use:

```bash
make docker-compose-dev-down
```

- To lint the project, use:

```bash
make lint
```

Please refer to the `Makefile` for more details on what each command does.