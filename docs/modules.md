# Documentation on Used Modules

This document describes the main modules and libraries used in the Kubernetes operator project.

## Table of Contents (Index)

### Main Modules
- **[Operator SDK](#operator-sdk)** - toolkit for creating Kubernetes operators
  - Description, version, and installation
  - CLI tools and scaffolding
  - Usage in the project
  - Advantages
- **[Go (Golang)](#go-golang)** - programming language
  - Description and features
  - Standard library (context, fmt, errors)
  - External dependencies (controller-runtime, k8s.io/*)

### External Go Dependencies
- **[sigs.k8s.io/controller-runtime](#sigsk8siocontroller-runtime)** - library for controllers
- **[k8s.io/api](#k8sioapi)** - Kubernetes API resource definitions
- **[k8s.io/apimachinery](#k8sioapimachinery)** - utilities for working with Kubernetes API
- **[k8s.io/client-go](#k8sioclient-go)** - client for interacting with Kubernetes API

### Project Management
- **[Dependency Management](#dependency-management)** - go.mod and commands
- **[Versioning](#versioning)** - SemVer and Go modules
- **[Development Tools](#development-tools)** - gofmt, goimports, golangci-lint, go test
- **[Update Recommendations](#update-recommendations)** - updating Operator SDK, Go, and dependencies

---

## Operator SDK

### Description
Operator SDK is a toolkit for creating, testing, and deploying Kubernetes operators. It provides high-level abstractions and tools to simplify operator development.

**Official Website**: https://sdk.operatorframework.io/  
**Installation Documentation**: https://sdk.operatorframework.io/docs/installation/  
**Repository**: https://github.com/operator-framework/operator-sdk

### Version
Recommended version: v1.42.0 or higher

### Installation

#### macOS (Homebrew)
```bash
brew install operator-sdk
```

#### Linux/Windows
```bash
# Set environment variables
export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')

# Download binary
export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.42.0
curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}

# Install
chmod +x operator-sdk_${OS}_${ARCH} && sudo mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk
```

### Main Components

#### 1. CLI Tools
- `operator-sdk init` - initialize a new project
- `operator-sdk create api` - create API and controller
- `operator-sdk generate` - generate code and manifests
- `operator-sdk build` - build operator image
- `operator-sdk run` - run operator locally

#### 2. Scaffolding
Operator SDK automatically generates:
- Project structure
- CRD definitions
- Controllers with basic logic
- RBAC manifests
- Kustomize configurations

#### 3. Integration with controller-runtime
Operator SDK uses controller-runtime for working with Kubernetes API.

### Usage in Project

**Detailed Guide**: See [Go Operator Creation Guide](operator-sdk-go-guide.md)

#### Project Initialization
```bash
cd operator
operator-sdk init --domain example.com --repo github.com/example/my-operator
```

#### API Creation
```bash
operator-sdk create api --group example --version v1 --kind MyResource
```

#### Code Generation
```bash
# Generate CRD manifests
operator-sdk generate kustomize manifests

# Generate controller code
operator-sdk generate
```

### Advantages
- Automatic boilerplate code generation
- Integration with OLM (Operator Lifecycle Manager)
- Support for various languages (Go, Ansible, Helm)
- Built-in testing tools
- Manifest generation for deployment

## Go (Golang)

### Description
Go is a compiled programming language with open source, developed by Google. Used as the main language for developing Kubernetes operators.

**Official Website**: https://go.dev/  
**Documentation**: https://go.dev/doc/

### Version
Required version: Go 1.22 or higher (project uses Go 1.23+)

**Note**: According to [Operator SDK official documentation](https://sdk.operatorframework.io/docs/building-operators/golang/installation/), the minimum Go version is 1.22. The project uses Go 1.23+ to access the latest language features.

### Installation
See official documentation: https://go.dev/doc/install

### Main Features for Operators

#### 1. Static Typing
- Ensures type safety at compile time
- Simplifies code refactoring and maintenance

#### 2. Simplicity and Readability
- Minimalistic syntax
- Explicit error handling
- No implicit conversions

#### 3. Performance
- Compilation to native code
- Efficient memory usage
- Fast execution time

#### 4. Concurrency
- Built-in goroutine support
- Channels for communication
- Excellent for asynchronous operations

### Standard Library Packages Used

#### context
- Operation lifecycle management
- Cancellation and timeouts
- Metadata passing

#### fmt
- String formatting
- Error handling with context

#### errors
- Error creation and wrapping
- Error type checking

### External Dependencies

#### sigs.k8s.io/controller-runtime
Main library for creating Kubernetes controllers.

**Version**: v0.18.0+

**Main Components**:
- `manager.Manager` - controller management
- `controller.Controller` - base controller
- `client.Client` - client for working with Kubernetes API
- `reconcile.Reconciler` - interface for reconciliation

**Usage**:
```go
import (
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller"
    "sigs.k8s.io/controller-runtime/pkg/manager"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
)
```

#### k8s.io/api
Kubernetes API resource definitions.

**Version**: v0.30.0+

**Usage**:
```go
import (
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
```

#### k8s.io/apimachinery
Utilities for working with Kubernetes API.

**Version**: v0.30.0+

**Main Components**:
- `runtime.Object` - interface for Kubernetes objects
- `meta/v1` - resource metadata
- `apis/meta/v1` - common API types

#### k8s.io/client-go
Client for interacting with Kubernetes API.

**Version**: v0.30.0+

**Usage**:
```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)
```

## Dependency Management

### go.mod
The `go.mod` file defines the module and its dependencies.

```go
module github.com/example/my-operator

go 1.23

require (
    sigs.k8s.io/controller-runtime v0.18.0
    k8s.io/api v0.30.0
    k8s.io/apimachinery v0.30.0
    k8s.io/client-go v0.30.0
)
```

### Commands for Dependency Management

```bash
# Add new dependency
go get package@version

# Update dependencies
go get -u ./...

# Clean unused dependencies
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

## Versioning

### Semantic Versioning
The project follows semantic versioning (SemVer):
- `MAJOR.MINOR.PATCH`
- Example: `v1.2.3`

### Go Modules and Versions
- Semantic versioning is used for dependencies
- Version tags must start with `v`
- `go.mod` automatically manages versions

## Development Tools

### gofmt
Code formatting:
```bash
gofmt -w .
```

### goimports
Import management:
```bash
goimports -w .
```

### golangci-lint
Static code analysis:
```bash
golangci-lint run
```

### go test
Run tests:
```bash
go test ./...
go test -v ./...
go test -cover ./...
```

## Update Recommendations

### Operator SDK
- Check changelog before updating
- Test updates in dev environment
- Follow migration guide for major updates

### Go
- Update to the latest stable version
- Check breaking changes in release notes
- Test dependency compatibility

### Dependencies
- Regularly update dependencies for security
- Use `go get -u` with caution
- Test after updating dependencies

---

**Note**: This document should be updated when module versions change or new dependencies are added.

