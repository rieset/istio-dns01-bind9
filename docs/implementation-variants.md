# Implementation Variants for DNS01 Challenge

## Overview

This operator supports two implementation variants for handling DNS01 challenges with multiple DNS servers:

1. **Variant 1: Cert-Manager Webhook Provider** (Implemented)
2. **Variant 2: Kubernetes Operator with CRD** (Planned)

## Variant 1: Cert-Manager Webhook Provider

### Description

A webhook provider for cert-manager that synchronously updates TXT records on multiple DNS servers using RFC2136 (nsupdate) protocol.

### Architecture

```
cert-manager → Webhook API → DNS01Solver → MultiServerDNS → RFC2136Client → DNS Servers
```

### Features

- ✅ Synchronous updates to multiple DNS servers
- ✅ TSIG authentication support
- ✅ Parallel updates for better performance
- ✅ Fault tolerance (minimum success threshold)
- ✅ Automatic cleanup after challenge completion

### Components

1. **DNS01Solver** (`internal/webhook/dns01_handler.go`)
   - Implements cert-manager webhook solver interface
   - Handles Present() and CleanUp() operations
   - Parses configuration and retrieves TSIG secrets

2. **MultiServerDNS** (`internal/webhook/multi_server.go`)
   - Manages updates across multiple DNS servers
   - Parallel execution for better performance
   - Fault tolerance with minimum success threshold

3. **RFC2136Client** (`internal/dns/rfc2136.go`)
   - RFC2136 protocol implementation
   - TSIG authentication
   - DNS update operations (add/delete TXT records)

### Usage

#### 1. Create TSIG Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tsig-secret
  namespace: cert-manager
type: Opaque
stringData:
  secret: "M730/BdKcT4VHgaISqojsqd/hdy1mdyPA8BQ824ASzo="
```

#### 2. Create ClusterIssuer or Issuer

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
              - "192.0.2.1"  # master
              - "192.0.2.2"  # secondary 1
              - "192.0.2.3"  # secondary 2
              - "192.0.2.4"  # secondary 3
            zone: "example.com"
            tsigKeyName: "acme-example-com"
            tsigAlgorithm: "hmac-sha256"
            tsigSecretName: "tsig-secret"
            tsigSecretKey: "secret"
            ttl: 60
```

#### 3. Create Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-cert
  namespace: default
spec:
  secretName: example-tls
  issuerRef:
    name: letsencrypt-dns01
    kind: ClusterIssuer
  dnsNames:
    - "app.example.com"
    - "*.example.com"
```

### Configuration Options

| Option | Required | Default | Description |
|--------|----------|---------|-------------|
| `servers` | Yes | - | List of DNS server IP addresses |
| `zone` | Yes | - | DNS zone name (e.g., "example.com") |
| `tsigKeyName` | Yes | - | TSIG key name |
| `tsigAlgorithm` | No | "hmac-sha256" | TSIG algorithm |
| `tsigSecretName` | Yes | - | Kubernetes Secret name containing TSIG secret |
| `tsigSecretKey` | No | "secret" | Key in Secret containing TSIG secret |
| `ttl` | No | 60 | TTL for TXT records |

### Deployment

The webhook solver can be deployed in two ways:

#### Option A: Separate Deployment (Recommended)

Deploy the webhook solver as a separate service:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns01-webhook-solver
  namespace: cert-manager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dns01-webhook-solver
  template:
    metadata:
      labels:
        app: dns01-webhook-solver
    spec:
      containers:
      - name: webhook
        image: example-registry.io/istio-dns01-bind9/operator-webhook:latest
        command: ["/manager", "webhook"]
        ports:
        - containerPort: 8089
          name: webhook
```

#### Option B: Integrated in Main Operator

The webhook solver can be integrated into the main operator by enabling it via flag.

### Advantages

- ✅ Direct integration with cert-manager
- ✅ No need for custom CRDs
- ✅ Automatic cleanup after challenge
- ✅ Works with existing cert-manager workflows
- ✅ Synchronous updates to all servers

### Limitations

- Requires cert-manager to be installed
- Webhook must be accessible from cert-manager
- Configuration is per Issuer/ClusterIssuer

## Variant 2: Kubernetes Operator with CRD (Planned)

### Description

A Kubernetes operator that manages DNS01 challenges through Custom Resources, providing more control and flexibility.

### Architecture

```
Certificate → DNS01Challenge CRD → Controller → MultiServerDNS → RFC2136Client → DNS Servers
```

### Features (Planned)

- Custom Resource Definition for DNS01 challenges
- Controller-based reconciliation
- Integration with cert-manager or standalone usage
- Status tracking and reporting
- Retry logic and error handling

### Usage (Planned)

```yaml
apiVersion: dns01.example.com/v1alpha1
kind: DNS01Challenge
metadata:
  name: example-challenge
  namespace: default
spec:
  domain: "app.example.com"
  token: "challenge-token"
  servers:
    - "192.0.2.1"
    - "192.0.2.2"
    - "192.0.2.3"
    - "192.0.2.4"
  zone: "example.com"
  tsigKeyName: "acme-example-com"
  tsigSecretRef:
    name: tsig-secret
    key: secret
```

### Advantages (Planned)

- Full control over challenge lifecycle
- Can work independently of cert-manager
- Better observability through CRD status
- More flexible for custom use cases

### Limitations (Planned)

- Requires creating and managing CRDs
- More complex setup
- Need to handle cleanup manually or via finalizers

## Comparison

| Feature | Variant 1 (Webhook) | Variant 2 (CRD) |
|---------|---------------------|-----------------|
| **Complexity** | Low | Medium |
| **Setup** | Simple | More complex |
| **Integration** | Direct with cert-manager | Flexible |
| **Control** | Limited to cert-manager | Full control |
| **Status Tracking** | Via cert-manager | Via CRD status |
| **Standalone** | No | Yes |
| **Implementation Status** | ✅ Implemented | ⏳ Planned |

## Recommendations

- **Use Variant 1** if you're using cert-manager and want simple, direct integration
- **Use Variant 2** if you need more control, want to work independently, or have custom requirements

## Implementation Status

- ✅ Variant 1: Fully implemented and ready for use
- ⏳ Variant 2: Planned for future implementation

## Documentation

- [Variant 1 Usage Guide](variant1-usage.md) - Detailed usage instructions for webhook provider
- [Variant 2 Design](variant2-design.md) - Design document for CRD-based operator (when implemented)

