# Istio DNS01 Bind9 Operator

Kubernetes operator for managing DNS01 challenges with Bind9 DNS servers in Istio environment.

## Description

This operator provides DNS01 challenge management for cert-manager, enabling synchronous DNS record updates across multiple Bind9 DNS servers using RFC2136 protocol with TSIG authentication. The operator works with multiple independent master DNS servers (without zone synchronization between them), updating each server directly. It supports two implementation variants:

1. **Cert-Manager Webhook Provider** (Implemented) - Direct integration with cert-manager
2. **Kubernetes Operator with CRD** (Planned) - Standalone operator with Custom Resources

## Features

- ✅ Synchronous DNS updates to multiple servers
- ✅ RFC2136 protocol support with TSIG authentication
- ✅ Parallel updates for better performance
- ✅ Fault tolerance with minimum success threshold
- ✅ Automatic cleanup after challenge completion
- ✅ Integration with cert-manager webhook framework

## Requirements

- Kubernetes cluster (version 1.19+)
- cert-manager installed (for Variant 1)
- DNS servers with Bind9 and TSIG keys configured
- **Important**: All DNS servers must be configured as master servers (without zone synchronization between them). The operator updates all servers directly via RFC2136 protocol.

## Project Structure

```
.
├── .cursorrules          # Rules for AI assistant
├── README.md             # This file
├── context.md            # Project context and state
├── docs/                 # Documentation
│   ├── implementation-variants.md
│   ├── variant1-usage.md
│   └── ...               # Other documentation
├── operator/             # Operator source code
│   ├── cmd/              # Entry points
│   ├── internal/         # Internal packages
│   │   ├── dns/          # RFC2136 DNS client
│   │   └── webhook/      # Cert-manager webhook solver
│   ├── config/           # Kubernetes manifests
│   └── ...               # Other operator files
└── ...                   # Other project files
```

## Quick Start

### Using Cert-Manager Webhook Provider (Variant 1)

1. **Create TSIG Secret**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tsig-secret
  namespace: cert-manager
type: Opaque
stringData:
  secret: "YOUR_TSIG_SECRET"
```

2. **Create ClusterIssuer**:
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns01
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          groupName: acme.example.com
          solverName: multi-dns
          config:
            servers:
              - "192.0.2.1"
              - "192.0.2.2"
            zone: "example.com"
            tsigKeyName: "acme-example-com"
            tsigAlgorithm: "hmac-sha256"
            tsigSecretName: "tsig-secret"
            tsigSecretKey: "secret"
```

3. **Create Certificate**:
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-cert
spec:
  secretName: example-tls
  issuerRef:
    name: letsencrypt-dns01
    kind: ClusterIssuer
  dnsNames:
    - "app.example.com"
```

For detailed instructions, see [Variant 1 Usage Guide](docs/variant1-usage.md).

## Documentation

Detailed documentation is located in the `docs/` directory:

### Implementation
- [Implementation Variants](docs/implementation-variants.md) - Overview of both implementation variants
- [Variant 1 Usage Guide](docs/variant1-usage.md) - Detailed usage instructions for webhook provider
- [Variant 2 Design](docs/variant2-design.md) - Design document for CRD-based operator (planned)

### Development
- [Best Practices](docs/best-practices.md) - Go development best practices
- [Function Rating System](docs/function-rating.md) - System for rating function quality and complexity
- [Operator SDK Guide](docs/operator-sdk-go-guide.md) - Guide to creating Go operators
- [Modules Documentation](docs/modules.md) - Documentation on used modules and dependencies
- [Linting Guide](docs/linting.md) - Code linting with golangci-lint
- [Logging Guide](docs/logging.md) - Logging system documentation
- [Glossary](docs/glossary.md) - Terms and definitions

## Deployment

### Deploy Webhook Solver (Variant 1)

Deploy the webhook solver to your Kubernetes cluster. See [Variant 1 Usage Guide](docs/variant1-usage.md) for detailed deployment instructions.

## Development

### Building the Operator

```bash
cd operator
make build
```

### Running Locally

```bash
cd operator
make run
```

### Testing

```bash
cd operator
make test
make lint
```

## License

This project is licensed under the MIT license.

## Authors

- **Albert Iblyaminov** - creator and main developer

