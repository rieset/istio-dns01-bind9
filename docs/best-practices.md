# Best Practices for Go Development

This document describes best practices used in the Kubernetes operator project.

## General Go Principles

### 1. Code Formatting
- Always use `gofmt` for code formatting
- Use `goimports` for automatic import management
- Configure IDE for automatic formatting on save

```go
// Good: proper formatting
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // ...
}

// Bad: improper formatting
func(r *Reconciler)Reconcile(ctx context.Context,req ctrl.Request)(ctrl.Result,error){
    // ...
}
```

### 2. Error Handling
- Always check errors explicitly
- Do not ignore errors using `_`
- Provide context in error messages
- Use `fmt.Errorf` with `%w` for error wrapping

```go
// Good: explicit error handling
if err != nil {
    return fmt.Errorf("failed to get pod: %w", err)
}

// Bad: ignoring errors
result, _ = client.Get(ctx, key, obj)
```

### 3. Using context.Context
- Always pass `context.Context` as the first parameter in functions that perform I/O operations
- Use context for operation cancellation and timeouts
- Do not store context in structures, pass it explicitly

```go
// Good: context as first parameter
func (c *Client) GetPod(ctx context.Context, name string) (*v1.Pod, error) {
    return c.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

// Bad: missing context
func (c *Client) GetPod(name string) (*v1.Pod, error) {
    return c.client.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
}
```

### 4. Naming
- Use short names for local variables
- Use long names for exported functions and types
- Follow Go naming conventions:
  - `Get`, `Set`, `Is`, `Has` for boolean methods
  - `New` for constructors
  - Interfaces may end with `-er` (e.g., `Reader`, `Writer`)

```go
// Good: proper naming
type PodMonitor struct {
    client client.Client
}

func (m *PodMonitor) IsPodReady(ctx context.Context, name string) (bool, error) {
    // ...
}

// Bad: improper naming
type pm struct {
    c client.Client
}

func (m *pm) check(ctx context.Context, n string) (bool, error) {
    // ...
}
```

### 5. Documentation
- **MANDATORY REQUIREMENT: All documentation and commits MUST be written in English**
  - All markdown files (README.md, docs/*.md) must be in English
  - All Go doc comments must be in English
  - All commit messages must be in English
  - All code comments should be in English
- All exported functions, types, and variables must have comments
- Comments should start with the entity name
- Use complete sentences
- **MANDATORY: Every function MUST have a FunctionRating comment - see Function Rating System section**

```go
// FunctionRating: 85/100
// - Complexity: LOW
// - Integrations: 1
// - External Risks: LOW
// - Unit Tests: YES
// - E2E Tests: YES
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: validateResource
// Purpose: Validates resource spec fields
// PodMonitor tracks the state of pods for the operator.
type PodMonitor struct {
    client client.Client
}

// FunctionRating: 90/100
// - Complexity: LOW
// - Integrations: 0
// - External Risks: LOW
// - Unit Tests: YES
// - E2E Tests: YES
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: Reconcile
// Purpose: Performs the reconciliation loop to bring the state to the desired state
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // ...
}
```

## Practices for Kubernetes Operators

### 1. Reconciliation Loop
- Always return `ctrl.Result{}` on successful completion
- Use `RequeueAfter` for periodic checks
- Use `Requeue: true` only for temporary errors
- Log all important actions

```go
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := ctrl.LoggerFrom(ctx)
    
    // Get resource
    instance := &examplev1.ExampleResource{}
    if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
        if apierrors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }
    
    // Execute logic
    if err := r.reconcileResources(ctx, instance); err != nil {
        log.Error(err, "failed to reconcile resources")
        return ctrl.Result{}, err
    }
    
    // Successful completion
    return ctrl.Result{}, nil
}
```

### 2. Working with Resources
- Always check resource existence before using
- Use `client.ObjectKey` for creating keys
- Handle `IsNotFound` errors separately
- Use patches instead of full updates when possible

```go
// Good: existence check
resource := &v1.Resource{}
key := client.ObjectKey{Namespace: namespace, Name: name}
if err := r.Get(ctx, key, resource); err != nil {
    if apierrors.IsNotFound(err) {
        // Resource does not exist, create new one
        return r.createResource(ctx, resource)
    }
    return err
}
```

### 3. Finalizers
- Use finalizers for resource cleanup
- Always remove finalizer after cleanup completion
- Handle resource deletion with finalizer in reconcile loop

```go
// Adding finalizer
if !controllerutil.ContainsFinalizer(instance, finalizerName) {
    controllerutil.AddFinalizer(instance, finalizerName)
    if err := r.Update(ctx, instance); err != nil {
        return err
    }
}

// Handling deletion
if !instance.GetDeletionTimestamp().IsZero() {
    // Perform cleanup
    if err := r.cleanup(ctx, instance); err != nil {
        return err
    }
    // Remove finalizer
    controllerutil.RemoveFinalizer(instance, finalizerName)
    return r.Update(ctx, instance)
}
```

### 4. Logging
- Use structured logging through `logr` (integrated in controller-runtime)
- Pass context through `context.Context`
- Use appropriate log levels:
  - `Error` - for errors requiring attention
  - `Info` - for important events
  - `Debug` - for debug information
- Include contextual information (resource names, namespace, host, domain, path)
- Avoid automatically added manifest fields - use only explicit fields

#### Logging Configuration

The project uses `zap` through `controller-runtime` with enhanced configuration:

```go
opts := zap.Options{
    Development: true,
    EncoderConfigOptions: []zap.EncoderConfigOption{
        func(config *zap.EncoderConfig) {
            // Color highlighting for log levels
            config.EncodeLevel = zap.CapitalColorLevelEncoder
            // ISO8601 time format
            config.EncodeTime = zap.ISO8601TimeEncoder
            // Caller format (file:line)
            config.EncodeCaller = zap.ShortCallerEncoder
        },
    },
}
```

#### Log Levels with Color Highlighting

- **ERROR** (red) - critical errors
- **INFO** (blue) - informational messages
- **DEBUG** (yellow) - debug information

#### Usage Examples

```go
// Get logger from context
logger := log.FromContext(ctx)

// Logging with contextual information
logger.Info("Resource detected",
    "resourceName", resource.Name,
    "resourceNamespace", resource.Namespace,
)

// Error logging
logger.Error(err, "failed to create resource",
    "resource", resource.Name,
    "namespace", namespace,
)
```

#### Log Typing

For better log type identification, use standard fields:
- `name`, `namespace` - basic fields for resources
- `resourceName`, `resourceNamespace` - resource information
- `podName`, `podNamespace` - Pod information
- `serviceName`, `serviceNamespace` - Service information

### 5. Testing
- Write unit tests for all business logic
- Use `envtest` for controller testing
- Mock external dependencies
- Use table-driven tests for multiple scenarios

### 6. Linting After Code Generation
- **MANDATORY REQUIREMENT: After any code generation, linting must be run**
- Code generation can introduce formatting issues or rule violations
- Always run `make lint` after code generation commands:
  - `make generate` - code generation (deepcopy, client)
  - `make manifests` - manifest generation (CRD, RBAC)
  - `operator-sdk generate` - operator code generation
  - `operator-sdk generate kustomize manifests` - Kustomize manifest generation
- Fix all linting errors before continuing development
- Use `make lint-fix` for automatic fixing of some issues

```bash
# Example workflow with code generation
make generate          # Code generation
make lint              # MANDATORY: linting check
make lint-fix          # Automatic fixes (if possible)
make manifests         # Manifest generation
make lint              # MANDATORY: re-check after manifests
```

```go
func TestReconcile(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*testing.T, client.Client)
        wantErr bool
    }{
        {
            name: "successful reconciliation",
            setup: func(t *testing.T, c client.Client) {
                // Setup test data
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Execute test
        })
    }
}
```

## Code Structure

### 1. File Size
- **CRITICAL REQUIREMENT: Files must be less than 500 lines**
- If a file exceeds 500 lines, refactoring is **MANDATORY**
- **MANDATORY REQUIREMENT: If a file exceeds 300 lines, you MUST create modules and extract functions into separate modules**
  - Files with 300+ lines require immediate modularization
  - Extract related functions into separate modules/packages
  - Each module should have a single, clear responsibility
  - Use interfaces to maintain clean boundaries between modules
- Split large files into smaller, focused modules
- Each file should have a single, clear responsibility
- When approaching the limit (450+ lines), start planning separation

```go
// Bad: file with 350+ lines (exceeds 300 line limit)
// controller.go (350 lines)
// - Reconcile
// - reconcileResources
// - findResources
// - updateStatus
// - cleanup
// - validateResource
// ... and many more functions

// Good: separation into multiple files with modularization
// controller.go (150 lines)
// - Reconcile
// - SetupWithManager
//
// internal/controller/reconciler.go (120 lines)
// - reconcileResources
// - findResources
// - updateStatus
//
// internal/controller/cleanup.go (80 lines)
// - cleanup
// - removeFinalizers
//
// internal/controller/validation.go (60 lines)
// - validateResource
// - validateSpec
```

### 1.1. Modularization When Exceeding 300 Lines

When a file exceeds 300 lines, you must:

1. **Identify logical function groups**
   - Group related functions by functionality
   - Identify dependencies between groups

2. **Create separate modules**
   - Create packages in `internal/` for business logic
   - Use clear package names (e.g., `internal/reconciler`, `internal/cleanup`)

3. **Use interfaces**
   - Define interfaces for interaction between modules
   - This simplifies testing and maintenance

4. **Modularization example:**

```go
// Before modularization: controller.go (350 lines)
type Reconciler struct {
    client client.Client
}

func (r *Reconciler) Reconcile(...) { /* 50 lines */ }
func (r *Reconciler) reconcileResources(...) { /* 80 lines */ }
func (r *Reconciler) findResources(...) { /* 60 lines */ }
func (r *Reconciler) updateStatus(...) { /* 40 lines */ }
func (r *Reconciler) cleanup(...) { /* 70 lines */ }
func (r *Reconciler) validateResource(...) { /* 50 lines */ }

// After modularization:

// controller.go (100 lines)
type Reconciler struct {
    client      client.Client
    reconciler  *internal.Reconciler
    cleanup     *internal.Cleanup
    validator   *internal.Validator
}

func (r *Reconciler) Reconcile(...) { /* uses modules */ }

// internal/reconciler/reconciler.go (120 lines)
package reconciler

type Reconciler struct {
    client client.Client
}

func (r *Reconciler) ReconcileResources(...) { /* ... */ }
func (r *Reconciler) FindResources(...) { /* ... */ }
func (r *Reconciler) UpdateStatus(...) { /* ... */ }

// internal/cleanup/cleanup.go (80 lines)
package cleanup

type Cleanup struct {
    client client.Client
}

func (c *Cleanup) Cleanup(...) { /* ... */ }

// internal/validator/validator.go (60 lines)
package validator

type Validator struct{}

func (v *Validator) ValidateResource(...) { /* ... */ }
```

### 2. Package Organization
- Group related code into packages
- Avoid circular dependencies
- Use interfaces for abstraction

```
controllers/
  ├── example_controller.go  # Main controller
  └── resource_monitor.go     # Resource monitoring logic

internal/
  ├── client/                # Clients for external services
  └── utils/                  # Utilities
```

### 3. Interfaces
- Use interfaces for testability
- Define interfaces at the point of use, not implementation
- Keep interfaces small and focused

```go
// Interface for resource monitoring
type ResourceMonitor interface {
    IsResourceReady(ctx context.Context, name string) (bool, error)
    GetResourceStatus(ctx context.Context, name string) (*ResourceStatus, error)
}
```

### 4. Dependency Handling
- Use dependency injection
- Pass dependencies through constructors
- Avoid global variables

```go
// Good: dependency injection
type Reconciler struct {
    client         client.Client
    resourceMonitor ResourceMonitor
}

func NewReconciler(client client.Client, resourceMonitor ResourceMonitor) *Reconciler {
    return &Reconciler{
        client:          client,
        resourceMonitor: resourceMonitor,
    }
}
```

## Performance

### 1. Caching
- Use caching for frequently requested data
- Configure proper indexes for search
- Clear cache when necessary

### 2. Batch Operations
- Group operations when possible
- Use `List` instead of multiple `Get` requests

### 3. Async Operations
- Use goroutines for independent operations
- Manage goroutine lifecycle through context

## Security

### 1. RBAC
- Follow the principle of least privilege
- Request only necessary permissions
- Document required permissions

### 2. Input Validation
- Always validate input data
- Use webhooks for CR validation
- Check access rights before operations

### 3. Secrets
- Never log secrets
- Use Kubernetes Secrets for storing sensitive data
- Secret rotation should be automated

## Function Rating System

### Overview
Every function in the codebase MUST have a rating comment that evaluates its quality, complexity, and risk level. This system helps identify and prioritize refactoring efforts, starting with the most critical functions.

### Rating Format
Each function must include a rating comment in the following format:

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

### Rating Criteria

#### 1. Complexity (0-25 points)
Evaluate based on cyclomatic complexity and code structure:
- **LOW (20-25 points)**: Simple logic, linear flow, < 5 branches
- **MEDIUM (15-19 points)**: Moderate logic, some conditionals, 5-10 branches
- **HIGH (10-14 points)**: Complex logic, many conditionals, 10-15 branches
- **CRITICAL (0-9 points)**: Very complex, deeply nested, > 15 branches

#### 2. Integrations (0-15 points)
Count external dependencies and integrations:
- **0 integrations (15 points)**: Pure function, no external calls
- **1-2 integrations (12-14 points)**: Minimal dependencies
- **3-5 integrations (8-11 points)**: Moderate dependencies
- **6+ integrations (0-7 points)**: Many dependencies

#### 3. External Risks (0-20 points)
Assess risks from external API calls, network operations, file I/O:
- **LOW (18-20 points)**: No external calls, or only safe internal calls
- **MEDIUM (12-17 points)**: External calls with error handling and retries
- **HIGH (6-11 points)**: External calls with limited error handling
- **CRITICAL (0-5 points)**: External calls without proper error handling, timeouts, or retries

#### 4. Unit Tests (0-15 points)
- **YES (15 points)**: Comprehensive unit tests with >80% coverage
- **PARTIAL (8-14 points)**: Some tests, but coverage <80%
- **NO (0-7 points)**: No unit tests or minimal coverage

#### 5. E2E Tests (0-10 points)
- **YES (10 points)**: Function is covered by E2E tests
- **PARTIAL (5-9 points)**: Partially covered by E2E tests
- **NO (0-4 points)**: No E2E test coverage

#### 6. Typing (0-15 points)
Evaluate type safety and explicit typing:
- **FULL (15 points)**: All parameters and returns are explicitly typed, no `interface{}`, proper generics
- **PARTIAL (8-14 points)**: Mostly typed, some `interface{}` usage
- **WEAK (0-7 points)**: Heavy use of `interface{}`, `any`, or untyped values

### Critical Issues (Automatic Score Reduction)

Functions with critical issues receive automatic penalties:

- **HARDCODE (-20 points)**: Contains hardcoded values (strings, numbers, URLs, etc.) that should be configurable
- **DEEP_NESTING (-15 points)**: Has nesting depth > 4 levels (if/for/switch/select)
- **BOTH (-35 points)**: Has both hardcode and deep nesting

### Score Calculation

Total Score = Complexity + Integrations + External Risks + Unit Tests + E2E Tests + Typing - Critical Issues Penalties

**Score Ranges:**
- **90-100**: Excellent - Well-tested, simple, safe function
- **75-89**: Good - Minor improvements possible
- **60-74**: Fair - Needs attention, consider refactoring
- **45-59**: Poor - Should be refactored soon
- **0-44**: Critical - Must be refactored immediately

### Rating Examples

#### Example 1: High-Rated Function

```go
// FunctionRating: 85/100
// - Complexity: LOW
// - Integrations: 1
// - External Risks: LOW
// - Unit Tests: YES
// - E2E Tests: YES
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: validateResource
// Purpose: Validates resource spec fields
func validateResource(spec *ResourceSpec) error {
    if spec.Name == "" {
        return fmt.Errorf("name is required")
    }
    if spec.Replicas < 0 {
        return fmt.Errorf("replicas must be non-negative")
    }
    return nil
}
```

#### Example 2: Low-Rated Function (Needs Refactoring)

```go
// FunctionRating: 35/100
// - Complexity: HIGH
// - Integrations: 5
// - External Risks: CRITICAL
// - Unit Tests: NO
// - E2E Tests: NO
// - Typing: WEAK
// - Critical Issues: BOTH
//
// Function: processResource
// Purpose: Processes resource with external API calls
func processResource(ctx context.Context, obj interface{}) error {
    // CRITICAL: Hardcoded URL
    url := "https://api.example.com/v1/process"
    
    // CRITICAL: Deep nesting (5 levels)
    if obj != nil {
        if data, ok := obj.(map[string]interface{}); ok {
            if name, ok := data["name"].(string); ok {
                if len(name) > 0 {
                    // CRITICAL: No error handling, no timeout
                    resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(name)))
                    if err != nil {
                        return err
                    }
                    defer resp.Body.Close()
                }
            }
        }
    }
    return nil
}
```

#### Example 3: Medium-Rated Function (Needs Improvement)

```go
// FunctionRating: 65/100
// - Complexity: MEDIUM
// - Integrations: 3
// - External Risks: MEDIUM
// - Unit Tests: PARTIAL
// - E2E Tests: NO
// - Typing: PARTIAL
// - Critical Issues: NONE
//
// Function: reconcileResources
// Purpose: Reconciles multiple resources
func (r *Reconciler) reconcileResources(ctx context.Context, instance *v1.MyResource) error {
    // Some complexity with multiple conditionals
    if instance.Spec.Enabled {
        if err := r.createService(ctx, instance); err != nil {
            return err
        }
        if instance.Spec.Expose {
            if err := r.createIngress(ctx, instance); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### Mandatory Requirements

1. **Every exported function MUST have a rating comment**
2. **Every function with complexity HIGH or CRITICAL MUST have a rating**
3. **Functions with critical issues (HARDCODE, DEEP_NESTING) MUST be refactored**
4. **Functions with score < 60 MUST be prioritized for refactoring**
5. **Functions with score < 45 MUST be refactored before merging**

### Refactoring Priority

When refactoring, prioritize functions in this order:
1. Functions with CRITICAL issues (HARDCODE, DEEP_NESTING)
2. Functions with score < 45
3. Functions with score 45-59
4. Functions with HIGH complexity
5. Functions with CRITICAL external risks
6. Functions without unit tests

### Rating Maintenance

- Update ratings when modifying functions
- Re-evaluate ratings during code reviews
- Track rating improvements in commit messages
- Maintain a list of low-rated functions for refactoring backlog

### How to Calculate Ratings

1. **Analyze the function**:
   - Count branches (if, for, switch, select)
   - Count external calls (API, network, file I/O)
   - Check for hardcoded values
   - Measure nesting depth
   - Check test coverage
   - Evaluate type safety

2. **Assign points** for each criterion

3. **Apply penalties** for critical issues

4. **Calculate total score**

5. **Add rating comment** above the function

### Best Practices for High Ratings

- Keep functions simple and focused
- Minimize external dependencies
- Always handle errors properly
- Write comprehensive tests
- Use explicit types
- Extract hardcoded values to configuration
- Reduce nesting depth (use early returns, extract functions)
- Add proper timeouts and retries for external calls

---

**Note**: This document should be updated as the project evolves and new practices emerge.

