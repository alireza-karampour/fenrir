# Buildah Go Library Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Basic Usage](#basic-usage)
4. [Core Concepts](#core-concepts)
5. [Image Building](#image-building)
6. [Storage Management](#storage-management)
7. [Advanced Features](#advanced-features)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [Examples](#examples)

## Introduction

Buildah is a tool that facilitates building OCI-compliant container images without requiring a Docker daemon or Dockerfile. When used as a Go library, it provides programmatic control over container image creation, allowing you to build images dynamically based on your application's requirements.

### Key Benefits

- **No Dockerfile Required**: Build images programmatically without writing Dockerfiles
- **OCI Compliant**: Creates standard OCI-compliant images
- **Rootless Support**: Can run without root privileges
- **Fine-grained Control**: Precise control over every aspect of image building
- **Integration**: Seamlessly integrates with Go applications

## Installation

### Prerequisites

- Go 1.19 or later
- Linux environment (Buildah is Linux-specific)
- Optional: fuse-overlayfs for rootless operation

### Dependencies

Add Buildah to your Go module:

```bash
go get github.com/containers/buildah
go get github.com/containers/storage
go get github.com/containers/image/v5
go get github.com/sirupsen/logrus
```

### Required Packages (Ubuntu/Debian)

```bash
sudo apt-get update
sudo apt-get install -y buildah fuse-overlayfs
```

### Required Packages (CentOS/RHEL/Fedora)

```bash
sudo dnf install -y buildah fuse-overlayfs
```

## Basic Usage

### Simple Image Creation

```go
package main

import (
    "context"
    "log"

    "github.com/containers/buildah"
    is "github.com/containers/image/storage"
    "github.com/containers/storage"
    "github.com/sirupsen/logrus"
)

func main() {
    // Initialize Buildah for rootless execution
    if buildah.InitReexec() {
        return
    }

    // Set up logging
    logger := logrus.New()
    logger.Level = logrus.DebugLevel

    // Create storage store
    storeOptions, err := storage.DefaultStoreOptions()
    if err != nil {
        log.Fatalf("Error getting store options: %v", err)
    }

    store, err := storage.GetStore(storeOptions)
    if err != nil {
        log.Fatalf("Error getting storage store: %v", err)
    }
    defer store.Shutdown(false)

    // Build image
    imageID, err := buildSimpleImage(store, logger)
    if err != nil {
        log.Fatalf("Error building image: %v", err)
    }

    log.Printf("Image created successfully: %s", imageID)
}

func buildSimpleImage(store storage.Store, logger *logrus.Logger) (string, error) {
    // Create builder options
    builderOptions := buildah.BuilderOptions{
        FromImage: "docker.io/library/alpine:latest",
        Logger:    logger,
    }

    // Create new builder
    builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
    if err != nil {
        return "", err
    }
    defer builder.Delete()

    // Configure container
    err = configureContainer(builder)
    if err != nil {
        return "", err
    }

    // Commit to image
    imageRef, err := is.Transport.ParseStoreReference(store, "my-app:latest")
    if err != nil {
        return "", err
    }

    imageID, _, _, err := builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    if err != nil {
        return "", err
    }

    return imageID, nil
}

func configureContainer(builder *buildah.Builder) error {
    // Install packages
    err := builder.Run([]string{"apk", "add", "--no-cache", "nginx"}, buildah.RunOptions{})
    if err != nil {
        return err
    }

    // Set environment variables
    builder.SetEnv("NGINX_HOST", "0.0.0.0")
    builder.SetEnv("NGINX_PORT", "80")

    // Set working directory
    builder.SetWorkingDir("/var/www/html")

    // Set entrypoint
    builder.SetEntrypoint([]string{"nginx"})
    builder.SetCmd([]string{"-g", "daemon off;"})

    // Expose port
    builder.SetPort("80/tcp")

    return nil
}
```

## Core Concepts

### Storage Store

The storage store manages container layers and images. It's the foundation of Buildah's operation.

```go
// Get default storage options
storeOptions, err := storage.DefaultStoreOptions()
if err != nil {
    return err
}

// Create store
store, err := storage.GetStore(storeOptions)
if err != nil {
    return err
}
defer store.Shutdown(false)
```

### Builder

The Builder is the core object for creating and manipulating container images.

```go
builderOptions := buildah.BuilderOptions{
    FromImage: "docker.io/library/alpine:latest", // Base image
    Logger:    logger,                            // Optional logger
}

builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
if err != nil {
    return err
}
defer builder.Delete() // Always clean up
```

### Image References

Images are referenced using OCI-compliant references:

```go
// Parse image reference
imageRef, err := is.Transport.ParseStoreReference(store, "my-app:latest")
if err != nil {
    return err
}

// Commit builder to image
imageID, _, _, err := builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
```

## Image Building

### Dockerfile Equivalent Operations

| Dockerfile Command | Buildah Method |
|-------------------|----------------|
| `FROM` | `BuilderOptions.FromImage` |
| `RUN` | `builder.Run()` |
| `COPY` | `builder.Add()` |
| `ADD` | `builder.Add()` |
| `ENV` | `builder.SetEnv()` |
| `WORKDIR` | `builder.SetWorkingDir()` |
| `USER` | `builder.SetUser()` |
| `ENTRYPOINT` | `builder.SetEntrypoint()` |
| `CMD` | `builder.SetCmd()` |
| `EXPOSE` | `builder.SetPort()` |

### Running Commands

```go
// Single command
err := builder.Run([]string{"apk", "add", "--no-cache", "nginx"}, buildah.RunOptions{})

// Multiple commands
commands := [][]string{
    {"apk", "update"},
    {"apk", "add", "--no-cache", "nginx", "nodejs"},
    {"npm", "install", "-g", "pm2"},
}

for _, cmd := range commands {
    err := builder.Run(cmd, buildah.RunOptions{})
    if err != nil {
        return err
    }
}
```

### Copying Files

```go
// Copy single file
err := builder.Add("/app/config.json", false, buildah.AddAndCopyOptions{}, "./config.json")

// Copy directory
err := builder.Add("/app/", false, buildah.AddAndCopyOptions{}, "./src")

// Copy with ownership
err := builder.Add("/app/", false, buildah.AddAndCopyOptions{
    Chown: "1000:1000",
}, "./src")
```

### Environment Configuration

```go
// Set single environment variable
builder.SetEnv("NODE_ENV", "production")

// Set multiple environment variables
envVars := map[string]string{
    "NODE_ENV": "production",
    "PORT": "3000",
    "LOG_LEVEL": "info",
}

for key, value := range envVars {
    builder.SetEnv(key, value)
}
```

### User and Permissions

```go
// Set user
builder.SetUser("1000:1000")

// Set working directory
builder.SetWorkingDir("/app")

// Set entrypoint and command
builder.SetEntrypoint([]string{"node"})
builder.SetCmd([]string{"server.js"})
```

## Storage Management

### Storage Locations

Buildah stores images in different locations based on user context:

- **Root user**: `/var/lib/containers/storage/`
- **Non-root user**: `$HOME/.local/share/containers/storage/`

### Custom Storage Configuration

```go
func createCustomStore(customPath string) (storage.Store, error) {
    storeOptions := storage.StoreOptions{
        GraphRoot:       customPath,                    // Where images are stored
        RunRoot:         customPath + "/run",           // Temporary runtime files
        GraphDriverName: "overlay",                     // Storage driver
        GraphDriverOptions: []string{
            "overlay.mount_program=/usr/bin/fuse-overlayfs",
        },
    }

    return storage.GetStore(storeOptions)
}
```

### Image Management

```go
func manageImages(store storage.Store) error {
    // List all images
    images, err := store.Images()
    if err != nil {
        return err
    }

    fmt.Printf("Found %d images:\n", len(images))
    for _, img := range images {
        fmt.Printf("- %s (ID: %s)\n", img.Names, img.ID)
    }

    // Find specific image
    for _, img := range images {
        for _, name := range img.Names {
            if name == "my-app:latest" {
                fmt.Printf("Found image: %s (ID: %s)\n", name, img.ID)
                
                // Get image details
                imgInfo, err := store.Image(img.ID)
                if err != nil {
                    return err
                }
                
                fmt.Printf("Created: %s\n", imgInfo.Created)
                fmt.Printf("Size: %d bytes\n", imgInfo.Size)
            }
        }
    }

    return nil
}
```

### Image Export

```go
func exportImage(store storage.Store, imageName string) error {
    // Export to tar file
    imageRef, err := is.Transport.ParseStoreReference(store, imageName)
    if err != nil {
        return err
    }

    // Save to tar file
    destRef, err := alltransports.ParseImageName("docker-archive:/tmp/my-image.tar:" + imageName)
    if err != nil {
        return err
    }

    _, err = copy.Image(context.TODO(), &copy.Options{}, destRef, imageRef)
    return err
}
```

## Advanced Features

### Multi-stage Builds

```go
func buildMultiStage(store storage.Store, logger *logrus.Logger) error {
    // Build stage
    buildOptions := buildah.BuilderOptions{
        FromImage: "golang:1.21-alpine",
        Logger:    logger,
    }
    
    buildBuilder, err := buildah.NewBuilder(context.TODO(), store, buildOptions)
    if err != nil {
        return err
    }
    defer buildBuilder.Delete()

    // Build the application
    err = buildBuilder.Run([]string{"go", "build", "-o", "app", "./cmd/main.go"}, buildah.RunOptions{})
    if err != nil {
        return err
    }

    // Runtime stage
    runtimeOptions := buildah.BuilderOptions{
        FromImage: "alpine:latest",
        Logger:    logger,
    }
    
    runtimeBuilder, err := buildah.NewBuilder(context.TODO(), store, runtimeOptions)
    if err != nil {
        return err
    }
    defer runtimeBuilder.Delete()

    // Copy binary from build stage
    err = runtimeBuilder.CopyFrom(buildBuilder, "/go/app", "/usr/local/bin/app", buildah.AddAndCopyOptions{})
    if err != nil {
        return err
    }

    runtimeBuilder.SetEntrypoint([]string{"/usr/local/bin/app"})
    
    // Commit runtime image
    imageRef, err := is.Transport.ParseStoreReference(store, "my-app:latest")
    if err != nil {
        return err
    }

    _, _, _, err = runtimeBuilder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    return err
}
```

### Dynamic Configuration

```go
type ImageConfig struct {
    BaseImage   string
    Packages    []string
    Environment map[string]string
    Files       map[string]string
    ImageName   string
    User        string
    WorkingDir  string
    Entrypoint  []string
    Cmd         []string
    Ports       []string
}

func buildWithConfig(store storage.Store, config ImageConfig, logger *logrus.Logger) error {
    builderOptions := buildah.BuilderOptions{
        FromImage: config.BaseImage,
        Logger:    logger,
    }

    builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
    if err != nil {
        return err
    }
    defer builder.Delete()

    // Install packages dynamically
    for _, pkg := range config.Packages {
        err = builder.Run([]string{"apk", "add", "--no-cache", pkg}, buildah.RunOptions{})
        if err != nil {
            return err
        }
    }

    // Set environment variables dynamically
    for key, value := range config.Environment {
        builder.SetEnv(key, value)
    }

    // Copy files dynamically
    for src, dest := range config.Files {
        err = builder.Add(dest, false, buildah.AddAndCopyOptions{}, src)
        if err != nil {
            return err
        }
    }

    // Set user and working directory
    if config.User != "" {
        builder.SetUser(config.User)
    }
    if config.WorkingDir != "" {
        builder.SetWorkingDir(config.WorkingDir)
    }

    // Set entrypoint and command
    if len(config.Entrypoint) > 0 {
        builder.SetEntrypoint(config.Entrypoint)
    }
    if len(config.Cmd) > 0 {
        builder.SetCmd(config.Cmd)
    }

    // Expose ports
    for _, port := range config.Ports {
        builder.SetPort(port)
    }

    return commitBuilder(builder, store, config.ImageName)
}

func commitBuilder(builder *buildah.Builder, store storage.Store, imageName string) error {
    imageRef, err := is.Transport.ParseStoreReference(store, imageName)
    if err != nil {
        return err
    }

    _, _, _, err = builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    return err
}
```

### Rootless Operation

```go
package main

import (
    "github.com/containers/buildah"
    "github.com/containers/storage/pkg/unshare"
)

func main() {
    // Initialize Buildah for rootless execution
    if buildah.InitReexec() {
        return
    }
    
    // Re-execute in user namespace if needed
    unshare.MaybeReexecUsingUserNamespace(false)
    
    // Your application code here...
}
```

## Best Practices

### Error Handling

```go
func safeBuild(store storage.Store) error {
    builder, err := buildah.NewBuilder(context.TODO(), store, buildah.BuilderOptions{
        FromImage: "alpine:latest",
    })
    if err != nil {
        return fmt.Errorf("failed to create builder: %w", err)
    }
    
    // Always clean up
    defer func() {
        if deleteErr := builder.Delete(); deleteErr != nil {
            log.Printf("Warning: failed to delete builder: %v", deleteErr)
        }
    }()

    // Your build logic here...
    return nil
}
```

### Resource Management

```go
func buildWithCleanup(store storage.Store) error {
    // Create builder
    builder, err := buildah.NewBuilder(context.TODO(), store, buildah.BuilderOptions{
        FromImage: "alpine:latest",
    })
    if err != nil {
        return err
    }
    
    // Ensure cleanup
    defer func() {
        if err := builder.Delete(); err != nil {
            log.Printf("Failed to delete builder: %v", err)
        }
    }()

    // Build operations...
    
    return nil
}
```

### Logging

```go
func setupLogging() *logrus.Logger {
    logger := logrus.New()
    
    // Set log level based on environment
    if os.Getenv("DEBUG") == "true" {
        logger.Level = logrus.DebugLevel
    } else {
        logger.Level = logrus.InfoLevel
    }
    
    // Set formatter
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    
    return logger
}
```

### Configuration Management

```go
type BuildConfig struct {
    BaseImage   string            `yaml:"base_image"`
    Packages    []string          `yaml:"packages"`
    Environment map[string]string `yaml:"environment"`
    Files       map[string]string `yaml:"files"`
    User        string            `yaml:"user"`
    WorkingDir  string            `yaml:"working_dir"`
    Entrypoint  []string          `yaml:"entrypoint"`
    Cmd         []string          `yaml:"cmd"`
    Ports       []string          `yaml:"ports"`
}

func loadConfig(configPath string) (*BuildConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var config BuildConfig
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## Troubleshooting

### Common Issues

#### Permission Denied

```bash
# Ensure user is in the containers group
sudo usermod -a -G containers $USER

# Log out and log back in
```

#### Storage Driver Issues

```go
// Use fuse-overlayfs for rootless operation
storeOptions := storage.StoreOptions{
    GraphDriverName: "overlay",
    GraphDriverOptions: []string{
        "overlay.mount_program=/usr/bin/fuse-overlayfs",
    },
}
```

#### Memory Issues

```go
// Limit memory usage during builds
builderOptions := buildah.BuilderOptions{
    FromImage: "alpine:latest",
    Logger:    logger,
    // Add memory limits if needed
}
```

### Debug Mode

```go
func enableDebugMode() {
    // Set debug environment variable
    os.Setenv("BUILDAH_ISOLATION", "chroot")
    
    // Enable debug logging
    logger := logrus.New()
    logger.Level = logrus.DebugLevel
}
```

### Storage Cleanup

```go
func cleanupStorage(store storage.Store) error {
    // List all images
    images, err := store.Images()
    if err != nil {
        return err
    }
    
    // Remove unused images
    for _, img := range images {
        if len(img.Names) == 0 { // Unnamed images
            err = store.DeleteImage(img.ID, true)
            if err != nil {
                log.Printf("Failed to delete image %s: %v", img.ID, err)
            }
        }
    }
    
    return nil
}
```

## Examples

### Web Application Image

```go
func buildWebApp(store storage.Store, logger *logrus.Logger) error {
    builderOptions := buildah.BuilderOptions{
        FromImage: "node:18-alpine",
        Logger:    logger,
    }

    builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
    if err != nil {
        return err
    }
    defer builder.Delete()

    // Install dependencies
    err = builder.Run([]string{"npm", "install", "-g", "pm2"}, buildah.RunOptions{})
    if err != nil {
        return err
    }

    // Copy application files
    err = builder.Add("/app/", false, buildah.AddAndCopyOptions{}, "./src")
    if err != nil {
        return err
    }

    // Install application dependencies
    err = builder.Run([]string{"npm", "install"}, buildah.RunOptions{})
    if err != nil {
        return err
    }

    // Configure environment
    builder.SetEnv("NODE_ENV", "production")
    builder.SetEnv("PORT", "3000")
    builder.SetWorkingDir("/app")
    builder.SetUser("1000:1000")
    builder.SetEntrypoint([]string{"pm2-runtime"})
    builder.SetCmd([]string{"start", "ecosystem.config.js"})
    builder.SetPort("3000/tcp")

    // Commit image
    imageRef, err := is.Transport.ParseStoreReference(store, "web-app:latest")
    if err != nil {
        return err
    }

    _, _, _, err = builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    return err
}
```

### Database Image

```go
func buildDatabase(store storage.Store, logger *logrus.Logger) error {
    builderOptions := buildah.BuilderOptions{
        FromImage: "postgres:15-alpine",
        Logger:    logger,
    }

    builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
    if err != nil {
        return err
    }
    defer builder.Delete()

    // Copy initialization scripts
    err = builder.Add("/docker-entrypoint-initdb.d/", false, buildah.AddAndCopyOptions{}, "./init-scripts")
    if err != nil {
        return err
    }

    // Copy configuration
    err = builder.Add("/etc/postgresql/postgresql.conf", false, buildah.AddAndCopyOptions{}, "./postgresql.conf")
    if err != nil {
        return err
    }

    // Set environment variables
    builder.SetEnv("POSTGRES_DB", "myapp")
    builder.SetEnv("POSTGRES_USER", "myuser")
    builder.SetEnv("POSTGRES_PASSWORD", "mypassword")
    builder.SetPort("5432/tcp")

    // Commit image
    imageRef, err := is.Transport.ParseStoreReference(store, "database:latest")
    if err != nil {
        return err
    }

    _, _, _, err = builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    return err
}
```

### Microservice Image Builder

```go
type MicroserviceConfig struct {
    Name        string
    BaseImage   string
    SourcePath  string
    BuildCmd    []string
    RuntimeCmd  []string
    Environment map[string]string
    Ports       []string
}

func buildMicroservice(store storage.Store, config MicroserviceConfig, logger *logrus.Logger) error {
    builderOptions := buildah.BuilderOptions{
        FromImage: config.BaseImage,
        Logger:    logger,
    }

    builder, err := buildah.NewBuilder(context.TODO(), store, builderOptions)
    if err != nil {
        return err
    }
    defer builder.Delete()

    // Copy source code
    err = builder.Add("/app/", false, buildah.AddAndCopyOptions{}, config.SourcePath)
    if err != nil {
        return err
    }

    // Build application
    for _, cmd := range config.BuildCmd {
        err = builder.Run(cmd, buildah.RunOptions{})
        if err != nil {
            return err
        }
    }

    // Set environment variables
    for key, value := range config.Environment {
        builder.SetEnv(key, value)
    }

    // Set runtime configuration
    builder.SetWorkingDir("/app")
    builder.SetUser("1000:1000")
    builder.SetEntrypoint(config.RuntimeCmd)

    // Expose ports
    for _, port := range config.Ports {
        builder.SetPort(port)
    }

    // Commit image
    imageRef, err := is.Transport.ParseStoreReference(store, config.Name+":latest")
    if err != nil {
        return err
    }

    _, _, _, err = builder.Commit(context.TODO(), imageRef, buildah.CommitOptions{})
    return err
}
```

## Conclusion

Buildah as a Go library provides powerful programmatic control over container image creation. It eliminates the need for Dockerfiles while offering fine-grained control over every aspect of the build process. With proper error handling, resource management, and following best practices, you can create robust container image building solutions that integrate seamlessly with your Go applications.

For more information, refer to:
- [Buildah Documentation](https://buildah.io/)
- [Buildah GitHub Repository](https://github.com/containers/buildah)
- [OCI Image Specification](https://github.com/opencontainers/image-spec)
