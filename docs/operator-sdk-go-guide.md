# Guide to Creating a Go Operator with Operator SDK

This document analyzes the process of creating a Kubernetes operator in Go using Operator SDK based on official documentation.

**Official Documentation**: https://sdk.operatorframework.io/docs/building-operators/golang/

## Prerequisites

According to the [official documentation](https://sdk.operatorframework.io/docs/building-operators/golang/installation/), the following are required to create a Go operator:

### Required Components
- **git** - version control system
- **Go version 1.22+** - programming language
  - ⚠️ **Important**: The project specifies Go 1.23+, which meets the requirements
- **Docker version 17.03+** - for building operator images
- **kubectl** - client for working with Kubernetes
- **Access to a Kubernetes cluster** of compatible version

### Operator SDK CLI
- Installed `operator-sdk` CLI tool
- See [installation guide](https://sdk.operatorframework.io/docs/installation/)

## Operator Creation Process

### Step 1: Project Initialization

Command to initialize a new project:

```bash
cd operator
operator-sdk init --domain <domain> --repo <repo>
```

**Parameters:**
- `--domain` - domain for grouping APIs (e.g., `example.com`)
- `--repo` - path to Go module repository (e.g., `github.com/example/my-operator`)

**What is created:**
- Basic project structure
- `go.mod` file with dependencies
- `Makefile` with build and deployment commands
- Kustomize configuration
- Basic test files

**Example:**
```bash
operator-sdk init --domain example.com --repo github.com/example/my-operator
```

### Step 2: Creating API and Controller

Command to create Custom Resource and controller:

```bash
operator-sdk create api --group <group> --version <version> --kind <Kind> --resource --controller
```

**Parameters:**
- `--group` - API group (e.g., `cache`, `cert`)
- `--version` - API version (e.g., `v1alpha1`, `v1`)
- `--kind` - resource type (e.g., `Memcached`, `MyResource`)
- `--resource` - create Custom Resource Definition
- `--controller` - create controller

**What is created:**
- CRD definition in `api/<version>/<kind>_types.go`
- Controller in `controllers/<kind>_controller.go`
- Basic tests in `controllers/<kind>_controller_test.go`
- RBAC manifests for resource access

**Example:**
```bash
operator-sdk create api --group example --version v1 --kind MyResource --resource --controller
```

### Step 3: Defining Custom Resource

After creating the API, define the Custom Resource structure in `api/v1/myresource_types.go`:

```go
// MyResourceSpec defines the desired state of MyResource
type MyResourceSpec struct {
    // Specification fields
    Name string `json:"name,omitempty"`
    // ...
}

// MyResourceStatus defines the observed state of MyResource
type MyResourceStatus struct {
    // Status fields
    Ready bool `json:"ready,omitempty"`
    // ...
}

// MyResource is a Custom Resource for the operator
type MyResource struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MyResourceSpec   `json:"spec,omitempty"`
    Status MyResourceStatus `json:"status,omitempty"`
}
```

**Important:**
- After changing types, run `make generate` to update code
- Run `make manifests` to generate CRD manifests

### Step 4: Implementing Controller Logic

The main operator logic is implemented in the controller's `Reconcile` method:

```go
func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := log.FromContext(ctx)
    
    // 1. Get Custom Resource
    instance := &examplev1.MyResource{}
    if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
        if apierrors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }
    
    // 2. Execute reconciliation logic
    // - Find resources
    // - Check their state
    // - Update status
    
    // 3. Return result
    return ctrl.Result{}, nil
}
```

**Key Points:**
- Use `context.Context` for operation cancellation
- Handle `IsNotFound` errors for deleted resources
- Logging through `log.FromContext(ctx)`
- Return `ctrl.Result{}` for successful completion

### Step 5: Setting Up Watch for Resources

To monitor resources, set up watch:

```go
func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.MyResource{}).
        Owns(&corev1.Pod{}).  // Watch for pods
        WithOptions(controller.Options{
            MaxConcurrentReconciles: 1,
        }).
        Complete(r)
}
```

**Using Predicates for Filtering:**

```go
import "sigs.k8s.io/controller-runtime/pkg/predicate"

func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    // Predicate for filtering resources
    resourcePredicate := predicate.NewPredicateFuncs(func(obj client.Object) bool {
        // Filtering logic
        return true
    })
    
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplev1.MyResource{}).
        Owns(&corev1.Pod{}, builder.WithPredicates(resourcePredicate)).
        Complete(r)
}
```

### Step 6: Testing

Operator SDK uses `envtest` from controller-runtime for testing without a real cluster:

```go
// controllers/suite_test.go
var (
    cfg     *rest.Config
    k8sClient client.Client
    testEnv  *envtest.Environment
)

func TestMain(m *testing.M) {
    testEnv = &envtest.Environment{
        CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
        ErrorIfCRDPathMissing: true,
    }
    
    cfg, _ = testEnv.Start()
    k8sClient, _ = client.New(cfg, client.Options{Scheme: scheme.Scheme})
    
    code := m.Run()
    testEnv.Stop()
    os.Exit(code)
}
```

**Controller Testing:**

```go
func TestMyResourceReconciler(t *testing.T) {
    // Setup test environment
    // Create test resources
    // Call Reconcile
    // Check results
}
```

### Step 7: Code and Manifest Generation

After changing types and controller:

```bash
# Generate code (deepcopy, client, etc.)
make generate

# Generate CRD manifests
make manifests

# Generate everything (code + manifests)
operator-sdk generate kustomize manifests
```

### Step 8: Local Run

For development and debugging:

```bash
# Run operator locally (requires cluster access)
make run

# Or with explicit namespace
make run NAMESPACE=default
```

**Requirements:**
- Configured `kubectl` with cluster access
- CRDs installed in cluster (`make install`)

### Step 9: Image Build

Build Docker image for operator:

```bash
# Build image
make docker-build IMG=<registry>/<name>:<tag>

# Example
make docker-build IMG=example-registry/my-operator:v0.1.0
```

**What happens:**
- Compilation of Go code into binary file
- Creation of Docker image with operator
- Use of base image from `Dockerfile`

### Step 10: Deployment

Deploy operator to cluster:

```bash
# Install CRDs
make install

# Deploy operator
make deploy IMG=<registry>/<name>:<tag>

# Example
make deploy IMG=example-registry/my-operator:v0.1.0
```

**What is created:**
- Namespace for operator
- ServiceAccount with RBAC permissions
- Deployment with operator
- Service (if required)

## Project Structure After Initialization

```
.
├── api/                    # API definitions
│   └── v1/
│       ├── myresource_types.go      # Type definitions
│       ├── myresource_types_test.go # Type tests
│       ├── groupversion_info.go      # Group version
│       └── zz_generated.deepcopy.go  # Auto-generated code
├── config/                  # Kustomize configuration
│   ├── crd/                # CRD manifests
│   ├── rbac/               # RBAC manifests
│   ├── manager/            # Operator deployment
│   └── samples/            # Custom Resource examples
├── controllers/            # Controllers
│   ├── myresource_controller.go     # Main controller
│   ├── myresource_controller_test.go # Controller tests
│   └── suite_test.go       # Test environment setup
├── main.go                 # Entry point
├── Makefile                # Build and deployment commands
├── go.mod                  # Go modules
├── go.sum                  # Dependency checksums
└── Dockerfile              # Docker operator image
```

## Key Makefile Commands

| Command | Description |
|---------|-------------|
| `make generate` | Generate code (deepcopy, client) |
| `make manifests` | Generate manifests (CRD, RBAC) |
| `make install` | Install CRDs to cluster |
| `make uninstall` | Remove CRDs from cluster |
| `make run` | Run operator locally |
| `make docker-build` | Build Docker image |
| `make docker-push` | Push image to registry |
| `make deploy` | Deploy operator |
| `make undeploy` | Remove operator |
| `make test` | Run tests |

## Best Practices from Documentation

### 1. Using Finalizers
For resource cleanup when deleting Custom Resource:

```go
const finalizerName = "myresource.example.com/finalizer"

if !controllerutil.ContainsFinalizer(instance, finalizerName) {
    controllerutil.AddFinalizer(instance, finalizerName)
    return r.Update(ctx, instance)
}

if !instance.GetDeletionTimestamp().IsZero() {
    // Perform cleanup
    controllerutil.RemoveFinalizer(instance, finalizerName)
    return r.Update(ctx, instance)
}
```

### 2. Status Update
Using `Status().Update()` to update status:

```go
instance.Status.Ready = true
instance.Status.Conditions = conditions
return r.Status().Update(ctx, instance)
```

### 3. Error Handling
Proper handling of different error types:

```go
if err != nil {
    if apierrors.IsNotFound(err) {
        // Resource not found - normal situation
        return ctrl.Result{}, nil
    }
    if apierrors.IsConflict(err) {
        // Version conflict - retry
        return ctrl.Result{Requeue: true}, nil
    }
    // Other errors
    return ctrl.Result{}, err
}
```

### 4. Logging
Structured logging with context:

```go
log := log.FromContext(ctx).WithValues(
    "myresource", req.NamespacedName,
    "namespace", req.Namespace,
)
log.Info("reconciling MyResource")
log.Error(err, "failed to reconcile")
```

## Documentation Links

- [Operator SDK Installation](https://sdk.operatorframework.io/docs/building-operators/golang/installation/)
- [Quickstart Tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/)
- [Testing with EnvTest](https://sdk.operatorframework.io/docs/building-operators/golang/testing/)
- [Advanced Topics](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/)
- [Controller Runtime Reference](https://sdk.operatorframework.io/docs/building-operators/golang/reference/)

---

**Note**: This document is based on Operator SDK v1.42.0 official documentation. When updating the version, check the current documentation.

