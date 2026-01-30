/*
Copyright 2026 Albert.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// FunctionRating: 85/100
// - Complexity: MEDIUM
// - Integrations: 3 (cert-manager, kubernetes client, webhook)
// - External Risks: MEDIUM (Kubernetes API, DNS operations)
// - Unit Tests: NO
// - E2E Tests: NO
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: DNS01Solver
// Purpose: Webhook solver for cert-manager DNS01 challenges with multi-server support

// DNS01Solver implements the cert-manager webhook solver interface
type DNS01Solver struct {
	client kubernetes.Interface
	logger *zap.Logger
}

// NewDNS01Solver creates a new DNS01 solver
func NewDNS01Solver(logger *zap.Logger) *DNS01Solver {
	return &DNS01Solver{
		logger: logger,
	}
}

// Name returns the name of the solver
func (s *DNS01Solver) Name() string {
	return "multi-dns"
}

// Present creates a TXT record for the DNS01 challenge
func (s *DNS01Solver) Present(ch *v1alpha1.ChallengeRequest) error {
	s.logger.Info("Presenting DNS01 challenge",
		zap.String("fqdn", ch.ResolvedFQDN),
		zap.String("key", ch.Key),
		zap.String("namespace", ch.ResourceNamespace),
	)

	// Parse configuration
	config, err := s.parseConfig(ch.Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Get TSIG secret from Kubernetes Secret
	tsigSecret, err := s.getTSIGSecret(ch.ResourceNamespace, config.TSIGSecretName, config.TSIGSecretKey)
	if err != nil {
		return fmt.Errorf("failed to get TSIG secret: %w", err)
	}

	// Create multi-server DNS manager
	dnsManager := NewMultiServerDNS(
		config.Servers,
		config.Zone,
		config.TSIGKeyName,
		config.TSIGAlgorithm,
		tsigSecret,
		s.logger,
	)

	// Add TXT record
	ctx := context.Background()
	if err := dnsManager.AddTXTRecord(ctx, ch.ResolvedFQDN, ch.Key, config.TTL); err != nil {
		return fmt.Errorf("failed to add TXT record: %w", err)
	}

	s.logger.Info("DNS01 challenge presented successfully",
		zap.String("fqdn", ch.ResolvedFQDN),
		zap.Int("servers", len(config.Servers)),
	)
	return nil
}

// CleanUp removes the TXT record after challenge completion
func (s *DNS01Solver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	s.logger.Info("Cleaning up DNS01 challenge",
		zap.String("fqdn", ch.ResolvedFQDN),
		zap.String("key", ch.Key),
		zap.String("namespace", ch.ResourceNamespace),
	)

	// Parse configuration
	config, err := s.parseConfig(ch.Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Get TSIG secret from Kubernetes Secret
	tsigSecret, err := s.getTSIGSecret(ch.ResourceNamespace, config.TSIGSecretName, config.TSIGSecretKey)
	if err != nil {
		return fmt.Errorf("failed to get TSIG secret: %w", err)
	}

	// Create multi-server DNS manager
	dnsManager := NewMultiServerDNS(
		config.Servers,
		config.Zone,
		config.TSIGKeyName,
		config.TSIGAlgorithm,
		tsigSecret,
		s.logger,
	)

	// Delete TXT record
	ctx := context.Background()
	if err := dnsManager.DeleteTXTRecord(ctx, ch.ResolvedFQDN); err != nil {
		return fmt.Errorf("failed to delete TXT record: %w", err)
	}

	s.logger.Info("DNS01 challenge cleaned up successfully",
		zap.String("fqdn", ch.ResolvedFQDN),
		zap.Int("servers", len(config.Servers)),
	)
	return nil
}

// Initialize initializes the solver with Kubernetes client
func (s *DNS01Solver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}
	s.client = cl
	return nil
}

// Config represents the webhook configuration
type Config struct {
	Servers        []string `json:"servers"`
	Zone           string   `json:"zone"`
	TSIGKeyName    string   `json:"tsigKeyName"`
	TSIGAlgorithm  string   `json:"tsigAlgorithm"`
	TSIGSecretName string   `json:"tsigSecretName"`
	TSIGSecretKey  string   `json:"tsigSecretKey"`
	TTL            int      `json:"ttl,omitempty"`
}

// parseConfig parses the webhook configuration
func (s *DNS01Solver) parseConfig(cfgJSON *acmev1.ACMESolverConfig) (*Config, error) {
	config := &Config{
		TTL:           60, // Default TTL
		TSIGAlgorithm: "hmac-sha256",
		TSIGSecretKey: "secret",
	}

	if cfgJSON == nil || len(cfgJSON.Raw) == 0 {
		return nil, fmt.Errorf("config is empty")
	}

	if err := json.Unmarshal(cfgJSON.Raw, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if len(config.Servers) == 0 {
		return nil, fmt.Errorf("servers list is required")
	}
	if config.Zone == "" {
		return nil, fmt.Errorf("zone is required")
	}
	if config.TSIGKeyName == "" {
		return nil, fmt.Errorf("tsigKeyName is required")
	}
	if config.TSIGSecretName == "" {
		return nil, fmt.Errorf("tsigSecretName is required")
	}

	return config, nil
}

// getTSIGSecret retrieves TSIG secret from Kubernetes Secret
func (s *DNS01Solver) getTSIGSecret(namespace, secretName, key string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("kubernetes client not initialized")
	}

	secret, err := s.client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s/%s: %w", namespace, secretName, err)
	}

	secretData, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s/%s", key, namespace, secretName)
	}

	return string(secretData), nil
}

// StartWebhookServer starts the webhook server
func StartWebhookServer(logger *zap.Logger) {
	groupName := "acme.example.com"
	solver := NewDNS01Solver(logger)

	cmd.RunWebhookServer(groupName,
		func(kubeClientConfig *rest.Config, stopCh <-chan struct{}) (v1alpha1.Solver, error) {
			if err := solver.Initialize(kubeClientConfig, stopCh); err != nil {
				return nil, err
			}
			return solver, nil
		},
	)
}

// HealthCheckHandler provides health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

