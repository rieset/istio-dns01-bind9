# Glossary of Terms

This file contains definitions of terms used in the project. **Important**: when introducing new terms, always add them to this file to synchronize understanding of entities.

## Core Terms

### Operator
A Kubernetes operator is a Kubernetes extension pattern that uses Custom Resources (CR) to manage applications and their components. The operator follows Kubernetes principles, particularly the concept of management through declarative APIs.

### Custom Resource (CR)
An extension of the Kubernetes API that is not part of the standard Kubernetes installation. Allows defining custom resource types.

### Custom Resource Definition (CRD)
A schema definition for a Custom Resource. CRD describes the data structure, validation, and versioning of a custom resource.

### Controller
A component of the operator that monitors resource changes and performs actions to bring the current state to the desired state (reconciliation loop).

### Reconciliation Loop
A process where the controller continuously checks the current state of resources and compares it with the desired state, making necessary changes.

### Reconcile
The process of bringing the current system state to the desired state described in the Custom Resource.

### Finalizer
A Kubernetes mechanism that allows a controller to perform cleanup before resource deletion. A resource with a finalizer cannot be deleted until the finalizer is removed.

### RBAC (Role-Based Access Control)
A role-based access control system in Kubernetes. The operator must have appropriate permissions to work with resources.

### Kustomize
A tool for customizing Kubernetes configurations. Used by Operator SDK for manifest generation.

### EnvTest
A test environment for controllers that provides a real Kubernetes API server without requiring a full cluster.

## Operator SDK Terms

### Operator SDK
A toolkit for creating, testing, and deploying Kubernetes operators. Supports Ansible, Helm, and Go operators.

### Controller Runtime
A library for creating Kubernetes controllers. The foundation for Go operators in Operator SDK.

### Manager
A controller-runtime component that manages the lifecycle of controllers and ensures their operation.

### Watch
A mechanism for tracking resource changes in Kubernetes. Controllers use watches to receive notifications about changes.

### Predicate
A filter for watch events that allows a controller to process only specific resource changes.

### Webhook
An HTTP server that processes requests from the Kubernetes API server. Used for resource validation and mutation.

## Kubernetes Resources

### Pod
The smallest deployment unit in Kubernetes. Represents one or more containers that work together and share resources.

### Service
An abstraction that defines a logical set of pods and access policy to them.

### Deployment
A resource that manages a replicated set of pods and ensures their updates.

### Namespace
A virtual cluster within a physical Kubernetes cluster. Used for resource isolation.

### Secret
An object containing sensitive data such as passwords, OAuth tokens, and SSH keys.

### ConfigMap
An object used to store configuration data as key-value pairs.

## Development Terms

### CRD (Custom Resource Definition)
A schema definition for a Custom Resource. Describes structure, validation, and versioning.

### Spec (Specification)
The desired state of a resource, defined by the user.

### Status
The observed state of a resource, updated by the operator.

### Owner Reference
A relationship between resources indicating the resource owner. Used for cascading deletion.

### Label
A label added to resources for organization and selection.

### Annotation
An annotation added to resources for storing metadata.

## Testing Terms

### Unit Test
A test that checks a single function or method in isolation.

### Integration Test
A test that checks interaction between system components.

### E2E Test (End-to-End Test)
A test that checks the complete system workflow from start to finish.

### Mock
An imitation of an object or function for testing.

### Test Fixture
A set of data or objects used for testing.

---

**Note**: This glossary should be updated as the project evolves. When adding new concepts, entities, or terms, always supplement this file.

