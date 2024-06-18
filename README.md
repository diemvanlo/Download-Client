# Project Guide

This guide will help you to run the `goload` project. The project is written in Go and uses a Makefile to automate various tasks such as generating code, building the project, and running the server.

## Prerequisites

- Go (version 1.16 or later)
- Docker and Docker Compose (for running the project in a containerized environment)
- golangci-lint (for linting Go code)
- buf (for generating API code)
- wire (for generating dependency injection code)

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