# Go Project Generator

A powerful and customizable Go project generator that helps you scaffold new Go services with a standardized structure and best practices.

## Features

- **Interactive CLI**: Guided setup through a user-friendly command-line interface
- **Modular Components**: Choose which components to include in your project
    - HTTP API with Gin
    - PostgreSQL database integration
    - Docker support with multi-stage builds
    - GitHub Actions CI/CD pipelines
- **Standardized Structure**: Follows Go project layout best practices
- **Database Migrations**: Built-in support for SQL migrations
- **Code Generation**: Automatic model generation from database schema
- **Git Integration**: Automatically initializes Git repository with GitHub remote

## Prerequisites

- Go 1.23 or higher
- Git

## Installation

### Using Go Install

```bash
go install github.com/neor-it/go-project-gen@latest
```

### Building from Source

```bash
git clone https://github.com/neor-it/go-project-gen.git
cd goprojectgen
go build -o goprojectgen main.go
```

### Using Docker

```bash
docker pull neor-it/goprojectgen:latest
```

Or build the Docker image:

```bash
docker build -t goprojectgen .
```

## Usage

### Direct Execution

```bash
goprojectgen
```

### Using Docker

```bash
# Using Docker volume to write output to current directory
docker run -it --rm -v $(pwd):/output goprojectgen

# Or with a specific Docker image
docker run -it --rm -v $(pwd):/output neor-it/goprojectgen:latest
```

When using Docker, the generated project will be created in your current directory, not inside the container.

### Using docker-compose

```bash
# Run the generator using docker-compose
docker-compose up
```

## Interactive Wizard

The generator will prompt you for the following information:

1. **GitHub username or organization**: Used for module path construction (e.g., `github.com/username/project-name`)
2. **Project name**: The name of your project and repository
3. **Components selection**: Choose which components to include:
    - HTTP server with Gin
    - PostgreSQL database
    - Docker support
    - CI/CD configuration

After confirming your choices, the generator will create the project structure with all the selected components.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

Created by the Go Project Generator Team