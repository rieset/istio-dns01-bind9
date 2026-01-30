# Variant 1: Cert-Manager Webhook Provider - Usage Guide

## Overview

This guide explains how to use the cert-manager webhook provider for DNS01 challenges with multiple DNS servers.

## Prerequisites

1. Kubernetes cluster (1.19+)
2. cert-manager installed
3. DNS servers configured with TSIG keys
4. Operator deployed with webhook solver
5. **Important**: All DNS servers must be configured as master servers (without zone synchronization between them). The operator updates all servers directly via RFC2136 protocol, so each server should have the zone configured as `type master` in Bind9 configuration.

## Step-by-Step Setup

### Step 1: Prepare TSIG Secret

Create a Kubernetes Secret containing your TSIG key:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tsig-secret
  namespace: cert-manager  # or your namespace
type: Opaque
stringData:
  secret: "M730/BdKcT4VHgaISqojsqd/hdy1mdyPA8BQ824ASzo="
```

**Important:** Replace the secret value with your actual TSIG secret.

### Step 2: Deploy Webhook Solver

Deploy the webhook solver as a service in your cluster:

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
      serviceAccountName: dns01-webhook-solver
      containers:
      - name: webhook
        image: example-registry.io/istio-dns01-bind9/operator-webhook:latest
        command: ["/manager", "webhook"]
        ports:
        - containerPort: 8089
          name: webhook
        env:
        - name: WEBHOOK_PORT
          value: "8089"
---
apiVersion: v1
kind: Service
metadata:
  name: dns01-webhook-solver
  namespace: cert-manager
spec:
  selector:
    app: dns01-webhook-solver
  ports:
  - port: 8089
    targetPort: 8089
    name: webhook
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dns01-webhook-solver
  namespace: cert-manager
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dns01-webhook-solver
  namespace: cert-manager
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dns01-webhook-solver
  namespace: cert-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dns01-webhook-solver
subjects:
- kind: ServiceAccount
  name: dns01-webhook-solver
  namespace: cert-manager
```

### Step 3: Register Webhook with cert-manager

Create a WebhookConfiguration to register the solver with cert-manager:

```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: dns01-webhook-solver
webhooks:
- name: dns01-webhook-solver.acme.example.com
  clientConfig:
    service:
      name: dns01-webhook-solver
      namespace: cert-manager
      path: "/"
  rules:
  - apiGroups: ["acme.cert-manager.io"]
    apiVersions: ["v1"]
    operations: ["CREATE", "UPDATE"]
    resources: ["challenges"]
  admissionReviewVersions: ["v1"]
  sideEffects: None
  failurePolicy: Fail
```

### Step 4: Create ClusterIssuer

Create a ClusterIssuer that uses the webhook solver:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns01
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@example.com
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          groupName: acme.example.com
          solverName: multi-dns
          config:
            servers:
              - "192.0.2.1"  # master DNS server 1
              - "192.0.2.2"  # master DNS server 2
              - "192.0.2.3"  # master DNS server 3
              - "192.0.2.4"  # master DNS server 4
            # Note: All servers must be configured as master (type master) 
            # without zone synchronization. The operator updates each server directly.
            zone: "example.com"
            tsigKeyName: "acme-example-com"
            tsigAlgorithm: "hmac-sha256"
            tsigSecretName: "tsig-secret"
            tsigSecretKey: "secret"
            ttl: 60
```

### Step 5: Create Certificate

Create a Certificate resource:

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

## Configuration Reference

### Webhook Config Structure

```json
{
  "servers": ["192.0.2.1", "192.0.2.2"],
  "zone": "example.com",
  "tsigKeyName": "acme-example-com",
  "tsigAlgorithm": "hmac-sha256",
  "tsigSecretName": "tsig-secret",
  "tsigSecretKey": "secret",
  "ttl": 60
}
```

### Field Descriptions

- **servers** (required): List of DNS server IP addresses to update. All servers must be configured as master servers (type master) in Bind9 without zone synchronization between them. The operator updates each server directly via RFC2136.
- **zone** (required): DNS zone name (e.g., "example.com")
- **tsigKeyName** (required): Name of the TSIG key configured on DNS servers
- **tsigAlgorithm** (optional): TSIG algorithm, default: "hmac-sha256"
- **tsigSecretName** (required): Kubernetes Secret name containing TSIG secret
- **tsigSecretKey** (optional): Key in Secret, default: "secret"
- **ttl** (optional): TTL for TXT records in seconds, default: 60

### DNS Server Configuration

**Important**: This operator is designed to work with multiple independent master DNS servers. Each server should be configured as follows:

- All servers must have the zone configured as `type master` in Bind9
- Zone synchronization (master-slave) between servers should **not** be enabled
- The operator updates each server directly via RFC2136 protocol
- Each server must have the same TSIG key configured for authentication

Example Bind9 configuration for each server:
```
zone "example.com" {
    type master;
    file "/etc/bind/zones/example.com.zone";
    allow-update { key "acme-example-com"; };
};
```

## Troubleshooting

### Check Webhook Solver Logs

```bash
kubectl logs -n cert-manager deployment/dns01-webhook-solver
```

### Check Certificate Status

```bash
kubectl describe certificate example-cert
kubectl describe challenge -n cert-manager
```

### Verify TXT Records

```bash
dig @192.0.2.1 TXT _acme-challenge.app.example.com +short
```

### Common Issues

1. **Webhook not called**: Check WebhookConfiguration and service
2. **TSIG authentication failed**: Verify TSIG secret and key name
3. **DNS update failed**: Check DNS server connectivity and zone configuration
4. **Some servers failed**: Check minimum success threshold (default: majority)

## Advanced Configuration

### Per-Domain Configuration

You can use different configurations for different domains:

```yaml
spec:
  acme:
    solvers:
    - selector:
        dnsZones:
          - "example.com"
      dns01:
        webhook:
          groupName: acme.example.com
          solverName: multi-dns
          config:
            servers: ["192.0.2.1", "192.0.2.2"]
            zone: "example.com"
            # ... other config
    - selector:
        dnsZones:
          - "another-domain.com"
      dns01:
        webhook:
          groupName: acme.example.com
          solverName: multi-dns
          config:
            servers: ["10.0.0.1", "10.0.0.2"]
            zone: "example.com"
            # ... other config
```

## Security Considerations

1. **TSIG Secrets**: Store TSIG secrets in Kubernetes Secrets, never in config
2. **RBAC**: Limit webhook solver permissions to only what's needed
3. **Network**: Ensure DNS servers are accessible from the cluster
4. **TLS**: Use TLS for webhook communication (cert-manager handles this)

## Performance

- Updates are performed in parallel across all servers
- Minimum success threshold: majority of servers (n/2 + 1)
- Timeout: 10 seconds per server
- Retry: Handled by cert-manager


