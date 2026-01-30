# Function Rating System

This document describes the function rating system used to evaluate code quality, complexity, and risk levels for all functions in the project.

## Purpose

The function rating system helps:
- Identify functions that need refactoring
- Prioritize refactoring efforts
- Track code quality improvements
- Ensure consistent code quality across the project
- Validate functions before merging

## Rating Format

Every function MUST include a rating comment in the following format:

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

## Rating Criteria

### 1. Complexity (0-25 points)

Evaluate based on cyclomatic complexity and code structure:

#### LOW (20-25 points)
- Simple logic, linear flow
- Less than 5 branches (if/for/switch/select)
- No nested conditionals beyond 2 levels
- Clear, straightforward implementation

**Example:**
```go
// FunctionRating: 22/25 (Complexity: LOW)
func validateName(name string) error {
    if name == "" {
        return fmt.Errorf("name is required")
    }
    if len(name) > 100 {
        return fmt.Errorf("name too long")
    }
    return nil
}
```

#### MEDIUM (15-19 points)
- Moderate logic with some conditionals
- 5-10 branches
- Some nested conditionals (2-3 levels)
- Manageable complexity

**Example:**
```go
// FunctionRating: 17/25 (Complexity: MEDIUM)
func processResource(resource *Resource) error {
    if resource == nil {
        return fmt.Errorf("resource is nil")
    }
    if resource.Spec.Enabled {
        if err := validateSpec(resource.Spec); err != nil {
            return err
        }
        if resource.Status.Ready {
            return r.updateResource(resource)
        }
    }
    return nil
}
```

#### HIGH (10-14 points)
- Complex logic with many conditionals
- 10-15 branches
- Deeply nested conditionals (3-4 levels)
- Multiple code paths

**Example:**
```go
// FunctionRating: 12/25 (Complexity: HIGH)
func reconcileResource(resource *Resource) error {
    if resource.Spec.Enabled {
        if resource.Status.Phase == PhasePending {
            if resource.Spec.Replicas > 0 {
                if err := r.createDeployment(resource); err != nil {
                    if apierrors.IsAlreadyExists(err) {
                        return r.updateDeployment(resource)
                    }
                    return err
                }
            }
        } else if resource.Status.Phase == PhaseRunning {
            // More nested logic...
        }
    }
    return nil
}
```

#### CRITICAL (0-9 points)
- Very complex logic
- More than 15 branches
- Deeply nested (4+ levels)
- Difficult to understand and maintain

**Example:**
```go
// FunctionRating: 5/25 (Complexity: CRITICAL)
func complexReconcile(resource *Resource) error {
    if condition1 {
        if condition2 {
            if condition3 {
                if condition4 {
                    if condition5 {
                        // Too many nested levels
                    }
                }
            }
        }
    }
    // Many more branches...
    return nil
}
```

### 2. Integrations (0-15 points)

Count external dependencies and integrations (API calls, database queries, file I/O, etc.):

#### 0 integrations (15 points)
- Pure function, no external calls
- Only uses parameters and local variables
- No side effects

**Example:**
```go
// FunctionRating: 15/15 (Integrations: 0)
func calculateTotal(items []Item) float64 {
    total := 0.0
    for _, item := range items {
        total += item.Price
    }
    return total
}
```

#### 1-2 integrations (12-14 points)
- Minimal dependencies
- One or two external calls
- Well-isolated

**Example:**
```go
// FunctionRating: 13/15 (Integrations: 2)
func getResource(ctx context.Context, name string) (*Resource, error) {
    resource := &Resource{}
    key := client.ObjectKey{Name: name}
    if err := r.client.Get(ctx, key, resource); err != nil {
        return nil, err
    }
    return resource, nil
}
```

#### 3-5 integrations (8-11 points)
- Moderate dependencies
- Multiple external calls
- Some coordination needed

**Example:**
```go
// FunctionRating: 10/15 (Integrations: 4)
func reconcileResource(ctx context.Context, resource *Resource) error {
    if err := r.client.Get(ctx, key, resource); err != nil {
        return err
    }
    if err := r.createService(ctx, resource); err != nil {
        return err
    }
    if err := r.createDeployment(ctx, resource); err != nil {
        return err
    }
    return r.updateStatus(ctx, resource)
}
```

#### 6+ integrations (0-7 points)
- Many dependencies
- Complex coordination
- High coupling

**Example:**
```go
// FunctionRating: 5/15 (Integrations: 8)
func complexReconcile(ctx context.Context, resource *Resource) error {
    // 8+ external calls
    r.client.Get(...)
    r.service.Create(...)
    r.deployment.Create(...)
    r.configMap.Create(...)
    r.secret.Create(...)
    r.serviceAccount.Create(...)
    r.role.Create(...)
    r.roleBinding.Create(...)
    // ...
}
```

### 3. External Risks (0-20 points)

Assess risks from external API calls, network operations, file I/O:

#### LOW (18-20 points)
- No external calls
- Only safe internal calls
- All operations are predictable

**Example:**
```go
// FunctionRating: 20/20 (External Risks: LOW)
func validateSpec(spec *Spec) error {
    if spec.Name == "" {
        return fmt.Errorf("name required")
    }
    return nil
}
```

#### MEDIUM (12-17 points)
- External calls with proper error handling
- Retries implemented
- Timeouts configured
- Context cancellation supported

**Example:**
```go
// FunctionRating: 15/20 (External Risks: MEDIUM)
func createResource(ctx context.Context, resource *Resource) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    if err := r.client.Create(ctx, resource); err != nil {
        if apierrors.IsAlreadyExists(err) {
            return nil // Already exists, OK
        }
        return fmt.Errorf("failed to create resource: %w", err)
    }
    return nil
}
```

#### HIGH (6-11 points)
- External calls with limited error handling
- No retries or timeouts
- Potential for failures

**Example:**
```go
// FunctionRating: 8/20 (External Risks: HIGH)
func callExternalAPI(url string, data []byte) error {
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
    if err != nil {
        return err // No retry, no timeout
    }
    defer resp.Body.Close()
    return nil
}
```

#### CRITICAL (0-5 points)
- External calls without proper error handling
- No timeouts or retries
- No context cancellation
- High risk of failures

**Example:**
```go
// FunctionRating: 2/20 (External Risks: CRITICAL)
func unsafeAPICall(url string) error {
    // No error handling, no timeout, no context
    resp, _ := http.Get(url)
    defer resp.Body.Close()
    return nil
}
```

### 4. Unit Tests (0-15 points)

#### YES (15 points)
- Comprehensive unit tests
- >80% code coverage
- Tests cover edge cases
- Tests are maintainable

**Example:**
```go
// FunctionRating: 15/15 (Unit Tests: YES)
// Function has comprehensive tests in *_test.go file
func validateResource(spec *Spec) error {
    // Implementation
}
```

#### PARTIAL (8-14 points)
- Some tests exist
- Coverage <80%
- Not all edge cases covered
- Tests may be outdated

**Example:**
```go
// FunctionRating: 10/15 (Unit Tests: PARTIAL)
// Function has some tests but coverage is ~60%
func processResource(resource *Resource) error {
    // Implementation
}
```

#### NO (0-7 points)
- No unit tests
- Minimal or no coverage
- Function is untested

**Example:**
```go
// FunctionRating: 0/15 (Unit Tests: NO)
// Function has no unit tests
func complexOperation(data interface{}) error {
    // Implementation
}
```

### 5. E2E Tests (0-10 points)

#### YES (10 points)
- Function is covered by E2E tests
- E2E tests verify end-to-end behavior
- Tests are reliable and maintained

**Example:**
```go
// FunctionRating: 10/10 (E2E Tests: YES)
// Function is tested in test/e2e/e2e_test.go
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Implementation
}
```

#### PARTIAL (5-9 points)
- Partially covered by E2E tests
- Some scenarios tested
- Not all paths covered

**Example:**
```go
// FunctionRating: 7/10 (E2E Tests: PARTIAL)
// Function has E2E tests but not all scenarios covered
func reconcileResource(ctx context.Context, resource *Resource) error {
    // Implementation
}
```

#### NO (0-4 points)
- No E2E test coverage
- Function behavior not verified end-to-end

**Example:**
```go
// FunctionRating: 0/10 (E2E Tests: NO)
// Function has no E2E tests
func processData(data []byte) error {
    // Implementation
}
```

### 6. Typing (0-15 points)

#### FULL (15 points)
- All parameters and returns explicitly typed
- No `interface{}` usage
- Proper use of generics (Go 1.18+)
- Type-safe throughout

**Example:**
```go
// FunctionRating: 15/15 (Typing: FULL)
func processResource(ctx context.Context, resource *v1.Resource) (*v1.ResourceStatus, error) {
    // All types explicit
    status := &v1.ResourceStatus{
        Ready: true,
    }
    return status, nil
}
```

#### PARTIAL (8-14 points)
- Mostly typed
- Some `interface{}` usage
- Some type assertions needed

**Example:**
```go
// FunctionRating: 12/15 (Typing: PARTIAL)
func processData(ctx context.Context, data interface{}) error {
    // Uses interface{} but has type assertions
    if str, ok := data.(string); ok {
        return processString(str)
    }
    return fmt.Errorf("unsupported type")
}
```

#### WEAK (0-7 points)
- Heavy use of `interface{}`, `any`
- Many type assertions
- Untyped values
- Type safety compromised

**Example:**
```go
// FunctionRating: 5/15 (Typing: WEAK)
func processAnything(data interface{}) (interface{}, error) {
    // Heavy use of interface{}, no type safety
    if m, ok := data.(map[string]interface{}); ok {
        if v, ok := m["value"].(interface{}); ok {
            // More type assertions...
        }
    }
    return nil, nil
}
```

## Critical Issues

Functions with critical issues receive automatic penalties:

### HARDCODE (-20 points)

Contains hardcoded values that should be configurable:
- Hardcoded strings (URLs, paths, names)
- Hardcoded numbers (timeouts, limits, ports)
- Hardcoded configuration values

**Examples of hardcode:**
```go
// BAD: Hardcoded URL
url := "https://api.example.com/v1/endpoint"

// BAD: Hardcoded timeout
timeout := 30 * time.Second

// BAD: Hardcoded port
port := 8080

// GOOD: From configuration
url := r.config.APIEndpoint
timeout := r.config.Timeout
port := r.config.Port
```

### DEEP_NESTING (-15 points)

Has nesting depth > 4 levels (if/for/switch/select):

**Example of deep nesting:**
```go
// BAD: 5 levels of nesting
if condition1 {
    if condition2 {
        if condition3 {
            if condition4 {
                if condition5 {
                    // Too deep!
                }
            }
        }
    }
}

// GOOD: Early returns reduce nesting
if !condition1 {
    return nil
}
if !condition2 {
    return nil
}
if !condition3 {
    return nil
}
// Continue with reduced nesting
```

### BOTH (-35 points)

Has both hardcode and deep nesting - most critical case.

## Score Calculation

**Total Score = Complexity + Integrations + External Risks + Unit Tests + E2E Tests + Typing - Critical Issues Penalties**

### Score Ranges

- **90-100**: Excellent
  - Well-tested, simple, safe function
  - No critical issues
  - Ready for production

- **75-89**: Good
  - Minor improvements possible
  - Generally well-written
  - May need some optimization

- **60-74**: Fair
  - Needs attention
  - Consider refactoring
  - Some improvements recommended

- **45-59**: Poor
  - Should be refactored soon
  - Has significant issues
  - Not ideal for production

- **0-44**: Critical
  - Must be refactored immediately
  - Has critical issues
  - High risk of bugs

## Complete Rating Examples

### Example 1: Excellent Function (95/100)

```go
// FunctionRating: 95/100
// - Complexity: LOW
// - Integrations: 1
// - External Risks: LOW
// - Unit Tests: YES
// - E2E Tests: YES
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: validateResourceSpec
// Purpose: Validates resource specification fields
func validateResourceSpec(spec *v1.ResourceSpec) error {
    if spec.Name == "" {
        return fmt.Errorf("name is required")
    }
    if spec.Replicas < 0 {
        return fmt.Errorf("replicas must be non-negative")
    }
    if spec.Replicas > 100 {
        return fmt.Errorf("replicas cannot exceed 100")
    }
    return nil
}
```

**Calculation:**
- Complexity: LOW = 22 points
- Integrations: 1 = 13 points
- External Risks: LOW = 20 points
- Unit Tests: YES = 15 points
- E2E Tests: YES = 10 points
- Typing: FULL = 15 points
- Critical Issues: NONE = 0 points
- **Total: 95/100**

### Example 2: Critical Function (25/100)

```go
// FunctionRating: 25/100
// - Complexity: CRITICAL
// - Integrations: 8
// - External Risks: CRITICAL
// - Unit Tests: NO
// - E2E Tests: NO
// - Typing: WEAK
// - Critical Issues: BOTH
//
// Function: processResource
// Purpose: Processes resource with multiple external calls
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
                    
                    // More external calls without proper handling
                    r.client.Create(ctx, obj)
                    r.service.Create(ctx, obj)
                    r.deployment.Create(ctx, obj)
                    // ... more calls
                }
            }
        }
    }
    return nil
}
```

**Calculation:**
- Complexity: CRITICAL = 5 points
- Integrations: 8 = 5 points
- External Risks: CRITICAL = 2 points
- Unit Tests: NO = 0 points
- E2E Tests: NO = 0 points
- Typing: WEAK = 5 points
- Critical Issues: BOTH = -35 points
- **Total: 25/100** (CRITICAL - must refactor)

### Example 3: Medium Function (68/100)

```go
// FunctionRating: 68/100
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
    if instance.Spec.Enabled {
        if err := r.createService(ctx, instance); err != nil {
            return err
        }
        if instance.Spec.Expose {
            if err := r.createIngress(ctx, instance); err != nil {
                return err
            }
        }
        return r.updateStatus(ctx, instance)
    }
    return nil
}
```

**Calculation:**
- Complexity: MEDIUM = 17 points
- Integrations: 3 = 10 points
- External Risks: MEDIUM = 15 points
- Unit Tests: PARTIAL = 10 points
- E2E Tests: NO = 0 points
- Typing: PARTIAL = 12 points
- Critical Issues: NONE = 0 points
- **Total: 68/100** (Fair - needs improvement)

## Mandatory Requirements

1. **Every exported function MUST have a rating comment**
2. **Every function with complexity HIGH or CRITICAL MUST have a rating**
3. **Functions with critical issues (HARDCODE, DEEP_NESTING) MUST be refactored**
4. **Functions with score < 60 MUST be prioritized for refactoring**
5. **Functions with score < 45 MUST be refactored before merging**

## Refactoring Priority

When refactoring, prioritize functions in this order:

1. **Functions with CRITICAL issues (HARDCODE, DEEP_NESTING)** - Fix immediately
2. **Functions with score < 45** - Must refactor before merging
3. **Functions with score 45-59** - Should refactor soon
4. **Functions with HIGH complexity** - Consider simplifying
5. **Functions with CRITICAL external risks** - Add proper error handling
6. **Functions without unit tests** - Add tests

## Rating Maintenance

- **Update ratings** when modifying functions
- **Re-evaluate ratings** during code reviews
- **Track rating improvements** in commit messages
- **Maintain a list** of low-rated functions for refactoring backlog
- **Set goals** for improving average function ratings

## Best Practices for High Ratings

1. **Keep functions simple**: Aim for LOW complexity
2. **Minimize dependencies**: Reduce integrations count
3. **Handle errors properly**: Add timeouts, retries, context cancellation
4. **Write comprehensive tests**: Aim for >80% coverage with unit and E2E tests
5. **Use explicit types**: Avoid `interface{}`, use generics
6. **Extract hardcoded values**: Move to configuration
7. **Reduce nesting**: Use early returns, extract functions
8. **Add proper error handling**: For all external calls

## Tools and Automation

Consider using tools to help calculate ratings:
- **gocyclo**: Measure cyclomatic complexity
- **go test -cover**: Check test coverage
- **golangci-lint**: Detect hardcoded values and nesting issues
- **Custom scripts**: Parse function ratings and generate reports

## See Also

- [Best Practices](./best-practices.md) - General development practices
- [.cursorrules](../.cursorrules) - Complete project rules including rating system

