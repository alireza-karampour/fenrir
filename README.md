# Fenrir

A powerful CLI tool for setting up end-to-end (E2E) test environments for Kubernetes-based applications. Fenrir automates the deployment and configuration of local Kubernetes clusters using Minikube, along with essential tools like kubectl and Helm.

## Features

- **Automated Kubernetes Setup**: Automatically downloads and configures Minikube for local Kubernetes development
- **Tool Management**: Downloads and manages kubectl and Helm with version verification
- **Test Suite Bootstrap**: Generates Ginkgo-based test suites with automatic cleanup
- **Chart Management**: Automatically installs Helm charts from the `charts/` directory
- **Image Loading**: Loads Docker images from local tar files into the cluster
- **Interactive Downloads**: Prompts for user confirmation before downloading tools
- **Checksum Verification**: Ensures downloaded binaries are authentic and uncorrupted

## Installation

### Prerequisites

- Go 1.24.5 or later
- Linux (amd64 architecture)
- Internet connection for downloading tools

### Build from Source

```bash
git clone https://github.com/alireza-karampour/fenrir.git
cd fenrir
go mod tidy
go build -o fenrir main.go
```

## Usage

### Basic Commands

```bash
# Initialize the test environment (downloads and configures all tools)
fenrir

# Bootstrap a test suite for the current package
fenrir --bootstrap

# Run tests (placeholder - functionality to be implemented)
fenrir test
```

### Command Options

- `--bootstrap, -b`: Generate a Ginkgo test suite for the current package
- `--help, -h`: Show help information

## Project Structure

```
fenrir/
├── bin/              # Downloaded binaries (minikube, kubectl, helm)
├── charts/           # Helm charts for deployment
├── cmd/              # Command definitions
│   ├── templates/    # Go templates for test generation
│   └── test/         # Test command implementation
├── images/           # Docker images (tar files) for cluster
├── pkg/              # Core packages
│   ├── cli/          # CLI utilities and subcommands
│   ├── task/         # Task execution framework
│   ├── utils/        # Utility functions
│   └── v1/           # Version-specific initialization
├── tars/             # Downloaded tar archives
├── main.go           # Application entry point
└── go.mod            # Go module definition
```

## Components

### Minikube Integration
- Downloads Minikube v1.36.0
- Automatically starts/stops clusters
- Enables MetalLB addon for load balancing
- Loads Docker images from `images/` directory

### Kubectl Integration
- Downloads kubectl v1.33.0
- Provides Kubernetes cluster interaction capabilities

### Helm Integration
- Downloads Helm v3.18.6
- Automatically installs charts from `charts/` directory
- Supports tar.gz extraction for binary distribution

### Test Framework
- **Bootstrap Feature**: Creates test suites using Go templates
- **Template System**: Uses `boot.gotmpl` template for generating Ginkgo test files
- **Test Structure**: Integrates with Ginkgo/Gomega testing framework
- **Auto-cleanup**: Automatically stops and deletes Minikube clusters after tests

## Configuration

### Environment Setup

Fenrir automatically manages the following tools:

| Tool | Version | Download Location |
|------|---------|-------------------|
| Minikube | v1.36.0 | `bin/minikube` |
| kubectl | v1.33.0 | `bin/kubectl` |
| Helm | v3.18.6 | `bin/helm` |

### Charts Directory

Place your Helm charts in the `charts/` directory. Fenrir will automatically install all charts found in subdirectories.

### Images Directory

Place Docker images as tar files in the `images/` directory. Fenrir will automatically load these images into the Minikube cluster.

## Development

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-ansi](https://codeberg.org/bit101/go-ansi) - Terminal colors
- [go-strcase](https://github.com/stoewer/go-strcase) - String case conversion

### Building

```bash
go mod tidy
go build -o fenrir main.go
```

### Testing

```bash
# Bootstrap a test suite
fenrir --bootstrap

# Run tests (when implemented)
go test ./...
```

## Workflow

1. **Initialization**: Run `fenrir` to download and configure all required tools
2. **Bootstrap**: Use `fenrir --bootstrap` to generate test suites for your packages
3. **Development**: Place charts in `charts/` and images in `images/`
4. **Testing**: Run your tests against the configured cluster
5. **Cleanup**: Tests automatically clean up the cluster after completion

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.

## Author

Created by [alireza-karampour](https://github.com/alireza-karampour)

---

**Note**: This tool is designed for Linux amd64 systems. Support for other platforms may be added in future versions.
