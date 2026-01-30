# Project Context: istio-dns01-bind9

This document contains the complete context of the `istio-dns01-bind9` Kubernetes operator project.

## General Information

**Project Type**: Kubernetes Operator (template, initialized via Operator SDK)

**Domain**: `istio-dns01-bind9.rieset.io`

**Repository**: `github.com/rieset/istio-dns01-bind9`

**Status**: Variant 1 (Cert-Manager Webhook Provider) implemented, Variant 2 (CRD-based) planned

**Author**: Albert Iblyaminov

**License**: Apache License 2.0

## Technical Specifications

- **Go Version**: 1.24.0+
- **Operator SDK**: v1.42.0+
- **Kubernetes**: 1.19+
- **Framework**: controller-runtime v0.21.0
- **Kubernetes API**: v0.33.0

## Project Structure

```
istio-dns01-bind9/
├── operator/              # Main operator code
│   ├── cmd/
│   │   ├── main.go        # Entry point (234 lines)
│   │   └── webhook/
│   │       └── main.go    # Webhook solver entry point
│   ├── internal/
│   │   ├── dns/
│   │   │   └── rfc2136.go # RFC2136 client implementation
│   │   └── webhook/
│   │       ├── dns01_handler.go  # Cert-manager webhook solver
│   │       └── multi_server.go   # Multi-server DNS manager
│   ├── config/            # Kustomize configurations
│   │   ├── crd/           # CRD definitions
│   │   ├── default/       # Default deployment config
│   │   ├── manager/       # Manager deployment
│   │   ├── rbac/          # RBAC configurations
│   │   ├── prometheus/    # Prometheus monitoring
│   │   ├── network-policy/# Network policies
│   │   ├── scorecard/     # Scorecard tests
│   │   └── manifests/     # OLM manifests
│   ├── test/
│   │   ├── e2e/           # E2E tests
│   │   └── utils/         # Test utilities
│   ├── hack/
│   │   └── boilerplate.go.txt  # License boilerplate
│   ├── go.mod             # Dependencies
│   ├── go.sum             # Dependency checksums
│   ├── Makefile           # Build and deployment commands
│   ├── Dockerfile         # Container image
│   ├── PROJECT            # Operator SDK project config
│   └── README.md          # Operator README
├── docs/                  # Documentation
│   ├── best-practices.md
│   ├── function-rating.md
│   ├── logging.md
│   ├── linting.md
│   ├── modules.md
│   ├── operator-sdk-go-guide.md
│   ├── glossary.md
│   ├── implementation-variants.md  # Implementation variants overview
│   └── variant1-usage.md          # Variant 1 usage guide
├── README.md              # Project overview
├── TEMPLATE.md            # Template information
├── context.md             # Project context (this file)
└── LICENSE                # MIT License
```

## Current State

### Implemented (Variant 1: Cert-Manager Webhook Provider)

- ✅ Basic operator structure (initialized via Operator SDK)
- ✅ Entry point `cmd/main.go` with metrics, webhooks, TLS support
- ✅ Cert-manager webhook solver (`internal/webhook/dns01_handler.go`)
- ✅ Multi-server DNS manager (`internal/webhook/multi_server.go`)
- ✅ RFC2136 client with TSIG support (`internal/dns/rfc2136.go`)
- ✅ Synchronous DNS updates to multiple servers
- ✅ Parallel updates with fault tolerance
- ✅ TSIG authentication support
- ✅ Automatic cleanup after challenge completion
- ✅ E2E tests (basic structure)
- ✅ Kustomize configurations (RBAC, Manager, Prometheus, Network Policy)
- ✅ Documentation on implementation variants and usage
- ✅ Makefile with build and deployment targets
- ✅ Dockerfile for container image

### Not Implemented

- ❌ Variant 2: Kubernetes Operator with CRD (planned)
- ❌ Unit tests for DNS operations
- ❌ Custom E2E test scenarios for DNS01 challenges
- ❌ Webhook deployment manifests (needs to be created)

## Development Rules

### 1. File Size Limits

- **Maximum**: 500 lines per file
- **Modularization required**: Files with 300+ lines MUST be modularized
- **Refactoring**: Files exceeding 500 lines require immediate refactoring

### 2. Function Rating System

Every function MUST have a `FunctionRating` comment:

```go
// FunctionRating: <SCORE>/100
// - Complexity: <LOW|MEDIUM|HIGH|CRITICAL>
// - Integrations: <COUNT>
// - External Risks: <LOW|MEDIUM|HIGH|CRITICAL>
// - Unit Tests: <YES|NO|PARTIAL>
// - E2E Tests: <YES|NO|PARTIAL>
// - Typing: <FULL|PARTIAL|WEAK>
// - Critical Issues: <NONE|HARDCODE|DEEP_NESTING|BOTH>
//
// Function: <function_name>
// Purpose: <brief_description>
```

**Rating Requirements:**
- Functions with score < 60: prioritize for refactoring
- Functions with score < 45: MUST refactor before merging
- Functions with CRITICAL issues: fix immediately

### 3. Linting

- **Tool**: `golangci-lint` v2.1.0
- **Configuration**: `.golangci.yml`
- **Mandatory**: Run linting after code generation
- **Auto-fix**: `make lint-fix`
- **Check**: `make lint`

### 4. Logging

- **Library**: `zap` via controller-runtime
- **Format**: Structured logging with color highlighting
- **Time format**: ISO8601
- **Standard fields**: `resourceName`, `namespace`, `host`, `path`
- **Levels**: ERROR (red), INFO (blue), DEBUG (yellow)

### 5. Documentation

- **Language**: All documentation and comments MUST be in English
- **Required**: Comments for all exported functions
- **Location**: Documentation in `docs/` directory
- **Format**: Markdown files

### 6. Code Generation

- **After generation**: MANDATORY linting check
- **Commands requiring linting**:
  - `make generate` - code generation (deepcopy, client)
  - `make manifests` - manifest generation (CRD, RBAC)
  - `operator-sdk generate` - operator code generation
  - `operator-sdk generate kustomize manifests` - Kustomize manifest generation

## Configuration

### Namespace

- **Default namespace**: `operator-system`
- **Name prefix**: `operator-`

### Service Account

- **Name**: `operator-controller-manager`
- **Namespace**: `operator-system`

### Metrics

- **Service name**: `operator-controller-manager-metrics-service`
- **Port**: 8443 (HTTPS)
- **Secure**: Enabled by default
- **Authentication**: RBAC-protected

### Leader Election

- **Enabled**: false (by default)
- **ID**: `5b98ccfb.istio-dns01-bind9.rieset.io`

### Health Probes

- **Address**: `:8081`
- **Endpoints**: `/healthz`, `/readyz`

## Dependencies

### Main Dependencies

```go
require (
    github.com/cert-manager/cert-manager/pkg/apis/acme/v1 v1.13.0
    github.com/miekg/dns v1.1.59
    github.com/onsi/ginkgo/v2 v2.22.0
    github.com/onsi/gomega v1.36.1
    k8s.io/apimachinery v0.33.0
    k8s.io/client-go v0.33.0
    sigs.k8s.io/controller-runtime v0.21.0
)
```

### Key Indirect Dependencies

- `go.uber.org/zap` - Structured logging
- `github.com/prometheus/client_golang` - Metrics
- `go.opentelemetry.io/otel` - Observability
- `sigs.k8s.io/kustomize` - Configuration management

## Makefile Targets

### Development

- `make build` - Build manager binary
- `make run` - Run controller locally
- `make generate` - Generate code (deepcopy, client)
- `make manifests` - Generate manifests (CRD, RBAC)
- `make fmt` - Format code
- `make vet` - Run go vet
- `make test` - Run unit tests
- `make lint` - Run linter
- `make lint-fix` - Run linter with auto-fix

### Build

- `make docker-build` - Build Docker image
- `make docker-push` - Push Docker image
- `make docker-buildx` - Build multi-platform image

### Deployment

- `make install` - Install CRDs
- `make uninstall` - Uninstall CRDs
- `make deploy` - Deploy controller
- `make undeploy` - Undeploy controller
- `make build-installer` - Generate install.yaml

### Testing

- `make test-e2e` - Run E2E tests (requires Kind)
- `make setup-test-e2e` - Setup Kind cluster for E2E tests
- `make cleanup-test-e2e` - Cleanup Kind cluster

### Dependencies

- `make kustomize` - Download kustomize
- `make controller-gen` - Download controller-gen
- `make envtest` - Download envtest
- `make golangci-lint` - Download golangci-lint
- `make operator-sdk` - Download operator-sdk

## Operator SDK Configuration

**PROJECT file:**
```yaml
domain: istio-dns01-bind9.rieset.io
layout:
- go.kubebuilder.io/v4
plugins:
  manifests.sdk.operatorframework.io/v2: {}
  scorecard.sdk.operatorframework.io/v2: {}
projectName: operator
repo: github.com/rieset/istio-dns01-bind9
version: "3"
```

## E2E Tests

**Location**: `operator/test/e2e/`

**Test Suite**: Ginkgo/Gomega

**Namespace**: `operator-system`

**Service Account**: `operator-controller-manager`

**Tests**:
- Manager pod running
- Metrics endpoint serving

## Documentation Files

### Best Practices (`docs/best-practices.md`)

- Code formatting
- Error handling
- Context usage
- Naming conventions
- Documentation requirements
- Reconciliation loop patterns
- Resource management
- Finalizers
- Logging configuration
- Testing practices
- Code structure and modularization
- Function rating system

### Function Rating System (`docs/function-rating.md`)

- Rating format and criteria
- Complexity evaluation
- Integration counting
- External risk assessment
- Test coverage requirements
- Typing evaluation
- Critical issues (hardcode, deep nesting)
- Score calculation
- Refactoring priorities

### Logging Guide (`docs/logging.md`)

- Zap configuration
- Log levels
- Standard fields
- Log filtering
- Production mode

### Linting Guide (`docs/linting.md`)

- golangci-lint setup
- Enabled linters
- CI/CD integration
- Best practices

### Operator SDK Guide (`docs/operator-sdk-go-guide.md`)

- Prerequisites
- Project initialization
- API creation
- Controller creation
- Testing

### Glossary (`docs/glossary.md`)

- Core Kubernetes terms
- Operator SDK terms
- Project-specific terminology

### Implementation Variants (`docs/implementation-variants.md`)

- Overview of both implementation variants
- Comparison and recommendations
- Architecture diagrams

### Variant 1 Usage Guide (`docs/variant1-usage.md`)

- Step-by-step setup instructions
- Configuration reference
- Troubleshooting guide
- Examples

## Implementation Variants

### Variant 1: Cert-Manager Webhook Provider (Implemented)

**Status**: ✅ Fully implemented

**Components**:
- `internal/webhook/dns01_handler.go` - Cert-manager webhook solver
- `internal/webhook/multi_server.go` - Multi-server DNS manager
- `internal/dns/rfc2136.go` - RFC2136 client with TSIG

**Features**:
- Synchronous DNS updates to multiple servers
- TSIG authentication
- Parallel updates with fault tolerance
- Automatic cleanup

**Usage**: See [docs/variant1-usage.md](docs/variant1-usage.md)

### Variant 2: Kubernetes Operator with CRD (Planned)

**Status**: ⏳ Planned for future implementation

**Components**: To be implemented

**Features** (planned):
- Custom Resource Definition for DNS01 challenges
- Controller-based reconciliation
- Standalone usage (independent of cert-manager)
- Status tracking and reporting

## Next Steps for Development

1. **Complete Variant 1 Deployment**
   - Create webhook deployment manifests
   - Add webhook service configuration
   - Test end-to-end with cert-manager

2. **Add Unit Tests**
   - Test RFC2136 client
   - Test multi-server DNS manager
   - Test webhook solver
   - Achieve >80% coverage

3. **Update E2E Tests**
   - Add DNS01 challenge test scenarios
   - Test multi-server updates
   - Test TSIG authentication

4. **Implement Variant 2** (Future)
   - Create DNS01Challenge CRD
   - Implement controller
   - Add reconciliation logic
   - Add status reporting

5. **Documentation**
   - Add deployment examples
   - Add troubleshooting guide
   - Add performance tuning guide

## Important Notes

1. **Project Independence**: This project is independent from other projects in the workspace. Do not use context from other projects.

2. **File Size**: Files must be less than 500 lines. Files with 300+ lines require modularization.

3. **Function Ratings**: All exported functions must have rating comments.

4. **Linting**: Always run linting after code generation.

5. **Documentation**: All documentation must be in English.

6. **Testing**: Write comprehensive unit and E2E tests.

7. **Error Handling**: Always handle errors explicitly with context.

8. **Context Usage**: Always pass `context.Context` as first parameter for I/O operations.

## Related Projects

This project is part of the `operators` workspace but is **independent**:
- `istio-gas/` - Istio HTTP01 Operator (separate project)
- `ai-set.kubernetes-operator/` - Another operator project (separate project)

**Important**: Do not mix context or code between these projects.

## Context File Maintenance

### IMPORTANT: Keep This File Updated

This `context.md` file is the **single source of truth** for project context. It MUST be kept up to date.

**When to Update:**
- After adding new API resources (CRDs)
- After implementing new controllers
- After adding/removing dependencies
- After changing project structure
- After updating configuration
- After adding new documentation
- After changing development workflow
- After updating technical specifications

**How to Update:**
1. Read the current `context.md` file
2. Identify sections that need updating
3. Update relevant sections with new information
4. Update the "Last Updated" date below
5. Ensure all information is accurate and complete

**AI Assistant Instructions:**
- ALWAYS read `context.md` before starting work on this project
- ALWAYS update `context.md` after making significant changes
- Use `context.md` as the primary reference for project state

---

**Last Updated**: 2026-01-29
**Context Version**: 2.0

