# Logging in Kubernetes Operator

This document describes the logging system used in the Kubernetes operator project.

## Overview

The project uses structured logging through the `zap` library (go.uber.org/zap), integrated into `controller-runtime`. Logging is configured for convenient development with color highlighting of levels and readable output format.

## Configuration

### Setup in `cmd/main.go`

```go
import (
    "go.uber.org/zap/zapcore"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

opts := zap.Options{
    Development: true,
    EncoderConfigOptions: []zap.EncoderConfigOption{
        func(config *zapcore.EncoderConfig) {
            // Color highlighting for log levels
            config.EncodeLevel = zapcore.CapitalColorLevelEncoder
            // ISO8601 time format
            config.EncodeTime = zapcore.ISO8601TimeEncoder
            // Caller format (file:line)
            config.EncodeCaller = zapcore.ShortCallerEncoder
        },
    },
}
```

### Command Line Flags

Zap provides the following flags for configuring logging:

- `--zap-devel` - enable development mode (color highlighting, readable format)
- `--zap-encoder` - encoder format (`json` or `console`)
- `--zap-log-level` - minimum log level (`debug`, `info`, `warn`, `error`)
- `--zap-stacktrace-level` - level for stack trace output (`error`, `panic`, `fatal`)

## Log Levels

### ERROR (red)
Used for critical errors requiring immediate attention.

```go
logger.Error(err, "failed to create resource",
    "name", resourceName,
    "namespace", namespace,
)
```

### INFO (blue)
Used for important events and informational messages.

```go
logger.Info("Resource detected",
    "resourceName", resource.Name,
    "resourceNamespace", resource.Namespace,
)
```

### DEBUG (yellow)
Used for debug information (enabled via `--zap-log-level=debug`).

```go
logger.V(1).Info("Processing resource",
    "resourceName", resource.Name,
    "state", currentState,
)
```

## Log Typing

For better log identification and filtering, standard fields are used:

### Fields for Resources

- `name`, `namespace` - basic fields for resources
- `resourceName`, `resourceNamespace` - resource information
- `podName`, `podNamespace` - Pod information
- `serviceName`, `serviceNamespace` - Service information
- `deploymentName`, `deploymentNamespace` - Deployment information

### Fields for Network Resources

- `host` - domain or host (e.g., `"app.example.com"`)
- `path` - URI path (e.g., `"/api/v1"`)
- `uri` - complete URI match information
- `destinationHost` - destination host
- `destinationPort` - destination port

## Log Examples

### Resource Detected

```
INFO    Resource detected    {"name": "example-resource", "namespace": "default"}
```

### Resource Created

```
INFO    Resource created    {"resourceName": "example-resource", "namespace": "default"}
```

### Error Creating Resource

```
ERROR   failed to create resource    {"error": "...", "resourceName": "example-resource", "namespace": "default"}
```

## Log Filtering

### By Resource Type

```bash
# Logs for specific resource
kubectl logs -n default deployment/my-operator | grep "Resource"
```

### By Resource Name

```bash
# Logs for specific resource
kubectl logs -n default deployment/my-operator | grep "resourceName.*example-resource"
```

### By Level

```bash
# Errors only
kubectl logs -n default deployment/my-operator | grep "ERROR"

# Informational messages only
kubectl logs -n default deployment/my-operator | grep "INFO"
```

## Production Mode

For production, JSON format is recommended:

```bash
# Run with JSON format
./manager --zap-encoder=json --zap-log-level=info
```

JSON format is convenient for parsing by monitoring and logging systems (ELK, Loki, etc.).

## Best Practices

1. **Always include contextual information**: resource names, namespace, identifiers
2. **Use appropriate levels**: Error for errors, Info for important events
3. **Avoid duplication**: do not log automatically added manifest fields
4. **Type logs**: use standard fields for better filtering
5. **Never log secrets**: never output sensitive data in logs

## See Also

- [Best Practices for Logging](./best-practices.md#4-logging)
- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Controller Runtime Logging](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/log)

