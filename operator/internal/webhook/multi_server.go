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
	"fmt"
	"sync"

	"github.com/rieset/istio-dns01-bind9/internal/dns"
	"go.uber.org/zap"
)

// FunctionRating: 80/100
// - Complexity: MEDIUM
// - Integrations: 1 (dns package)
// - External Risks: MEDIUM (multiple network operations)
// - Unit Tests: NO
// - E2E Tests: NO
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: MultiServerDNS
// Purpose: Manages DNS updates across multiple DNS servers synchronously

// MultiServerDNS handles DNS updates on multiple servers
type MultiServerDNS struct {
	servers  []string
	zone     string
	tsigKey  string
	tsigAlg  string
	tsigSec  string
	logger   *zap.Logger
	minSuccess int // Minimum number of successful updates required
}

// NewMultiServerDNS creates a new multi-server DNS manager
func NewMultiServerDNS(servers []string, zone, tsigKey, tsigAlg, tsigSec string, logger *zap.Logger) *MultiServerDNS {
	minSuccess := len(servers)/2 + 1 // At least half + 1 must succeed
	if minSuccess < 1 {
		minSuccess = 1
	}
	return &MultiServerDNS{
		servers:    servers,
		zone:       zone,
		tsigKey:    tsigKey,
		tsigAlg:    tsigAlg,
		tsigSec:    tsigSec,
		logger:     logger,
		minSuccess: minSuccess,
	}
}

// AddTXTRecord adds a TXT record to all configured DNS servers synchronously
func (m *MultiServerDNS) AddTXTRecord(ctx context.Context, fqdn, value string, ttl int) error {
	m.logger.Info("Adding TXT record to multiple servers",
		zap.String("fqdn", fqdn),
		zap.String("value", value),
		zap.Strings("servers", m.servers),
	)

	var wg sync.WaitGroup
	errChan := make(chan error, len(m.servers))
	successCount := 0
	var mu sync.Mutex

	for _, server := range m.servers {
		wg.Add(1)
		go func(srv string) {
			defer wg.Done()
			client := dns.NewRFC2136Client(srv, m.zone, m.tsigKey, m.tsigAlg, m.tsigSec, m.logger)
			if err := client.AddTXTRecord(ctx, fqdn, value, ttl); err != nil {
				m.logger.Error("Failed to add TXT record on server",
					zap.String("server", srv),
					zap.String("fqdn", fqdn),
					zap.Error(err),
				)
				errChan <- fmt.Errorf("server %s: %w", srv, err)
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
				m.logger.Info("Successfully added TXT record on server",
					zap.String("server", srv),
					zap.String("fqdn", fqdn),
				)
			}
		}(server)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// Check if we have enough successful updates
	if successCount < m.minSuccess {
		m.logger.Error("Insufficient successful updates",
			zap.Int("success_count", successCount),
			zap.Int("min_required", m.minSuccess),
			zap.Int("total_servers", len(m.servers)),
			zap.Int("errors", len(errors)),
		)
		return fmt.Errorf("only %d/%d servers updated successfully (minimum %d required): %v",
			successCount, len(m.servers), m.minSuccess, errors)
	}

	if len(errors) > 0 {
		m.logger.Warn("Some servers failed, but minimum success threshold met",
			zap.Int("success_count", successCount),
			zap.Int("error_count", len(errors)),
		)
	}

	m.logger.Info("TXT record added successfully to multiple servers",
		zap.String("fqdn", fqdn),
		zap.Int("success_count", successCount),
		zap.Int("total_servers", len(m.servers)),
	)
	return nil
}

// DeleteTXTRecord deletes a TXT record from all configured DNS servers synchronously
func (m *MultiServerDNS) DeleteTXTRecord(ctx context.Context, fqdn string) error {
	m.logger.Info("Deleting TXT record from multiple servers",
		zap.String("fqdn", fqdn),
		zap.Strings("servers", m.servers),
	)

	var wg sync.WaitGroup
	errChan := make(chan error, len(m.servers))
	successCount := 0
	var mu sync.Mutex

	for _, server := range m.servers {
		wg.Add(1)
		go func(srv string) {
			defer wg.Done()
			client := dns.NewRFC2136Client(srv, m.zone, m.tsigKey, m.tsigAlg, m.tsigSec, m.logger)
			if err := client.DeleteTXTRecord(ctx, fqdn); err != nil {
				m.logger.Error("Failed to delete TXT record on server",
					zap.String("server", srv),
					zap.String("fqdn", fqdn),
					zap.Error(err),
				)
				errChan <- fmt.Errorf("server %s: %w", srv, err)
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
				m.logger.Info("Successfully deleted TXT record on server",
					zap.String("server", srv),
					zap.String("fqdn", fqdn),
				)
			}
		}(server)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	// For deletion, we're more lenient - at least one server should succeed
	if successCount == 0 {
		m.logger.Error("Failed to delete TXT record on all servers",
			zap.Int("total_servers", len(m.servers)),
			zap.Int("errors", len(errors)),
		)
		return fmt.Errorf("failed to delete on all servers: %v", errors)
	}

	if len(errors) > 0 {
		m.logger.Warn("Some servers failed during deletion",
			zap.Int("success_count", successCount),
			zap.Int("error_count", len(errors)),
		)
	}

	m.logger.Info("TXT record deleted successfully from multiple servers",
		zap.String("fqdn", fqdn),
		zap.Int("success_count", successCount),
		zap.Int("total_servers", len(m.servers)),
	)
	return nil
}

