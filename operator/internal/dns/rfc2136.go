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

package dns

import (
	"context"
	"fmt"
	"time"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// FunctionRating: 75/100
// - Complexity: MEDIUM
// - Integrations: 2 (dns library, logging)
// - External Risks: MEDIUM (network operations, DNS server availability)
// - Unit Tests: NO
// - E2E Tests: NO
// - Typing: FULL
// - Critical Issues: NONE
//
// Function: RFC2136Client
// Purpose: Client for RFC2136 dynamic DNS updates using TSIG authentication

// RFC2136Client handles DNS updates via RFC2136 protocol
type RFC2136Client struct {
	server   string
	zone     string
	tsigKey  string
	tsigAlg  string
	tsigSec  string
	logger   *zap.Logger
	timeout  time.Duration
}

// NewRFC2136Client creates a new RFC2136 client
func NewRFC2136Client(server, zone, tsigKey, tsigAlg, tsigSec string, logger *zap.Logger) *RFC2136Client {
	return &RFC2136Client{
		server:  server,
		zone:    zone,
		tsigKey: tsigKey,
		tsigAlg: tsigAlg,
		tsigSec: tsigSec,
		logger:  logger,
		timeout: 10 * time.Second,
	}
}

// AddTXTRecord adds a TXT record to the DNS zone
func (c *RFC2136Client) AddTXTRecord(ctx context.Context, fqdn, value string, ttl int) error {
	c.logger.Info("Adding TXT record",
		zap.String("fqdn", fqdn),
		zap.String("value", value),
		zap.String("server", c.server),
		zap.String("zone", c.zone),
	)

	// Create DNS message
	msg := new(dns.Msg)
	msg.SetUpdate(dns.Fqdn(c.zone))

	// Create TXT record
	rr := new(dns.TXT)
	rr.Hdr = dns.RR_Header{
		Name:   dns.Fqdn(fqdn),
		Rrtype: dns.TypeTXT,
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}
	rr.Txt = []string{value}

	// Add record to message
	msg.Insert([]dns.RR{rr})

	// Add TSIG signature
	msg.SetTsig(c.tsigKey, c.tsigAlg, 300, time.Now().Unix())

	// Send update
	client := new(dns.Client)
	client.Timeout = c.timeout
	client.TsigSecret = map[string]string{c.tsigKey: c.tsigSec}

	reply, _, err := client.ExchangeContext(ctx, msg, c.server+":53")
	if err != nil {
		c.logger.Error("Failed to send DNS update",
			zap.String("fqdn", fqdn),
			zap.String("server", c.server),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send DNS update to %s: %w", c.server, err)
	}

	if reply.Rcode != dns.RcodeSuccess {
		c.logger.Error("DNS update failed",
			zap.String("fqdn", fqdn),
			zap.String("server", c.server),
			zap.Int("rcode", reply.Rcode),
			zap.String("rcode_name", dns.RcodeToString[reply.Rcode]),
		)
		return fmt.Errorf("DNS update failed: %s (rcode: %d)", dns.RcodeToString[reply.Rcode], reply.Rcode)
	}

	c.logger.Info("TXT record added successfully",
		zap.String("fqdn", fqdn),
		zap.String("server", c.server),
	)
	return nil
}

// DeleteTXTRecord deletes a TXT record from the DNS zone
func (c *RFC2136Client) DeleteTXTRecord(ctx context.Context, fqdn string) error {
	c.logger.Info("Deleting TXT record",
		zap.String("fqdn", fqdn),
		zap.String("server", c.server),
		zap.String("zone", c.zone),
	)

	// Create DNS message
	msg := new(dns.Msg)
	msg.SetUpdate(dns.Fqdn(c.zone))

	// Create RR to delete
	rr := new(dns.TXT)
	rr.Hdr = dns.RR_Header{
		Name:   dns.Fqdn(fqdn),
		Rrtype: dns.TypeTXT,
		Class:  dns.ClassINET,
	}

	// Remove record
	msg.RemoveRRset([]dns.RR{rr})

	// Add TSIG signature
	msg.SetTsig(c.tsigKey, c.tsigAlg, 300, time.Now().Unix())

	// Send update
	client := new(dns.Client)
	client.Timeout = c.timeout
	client.TsigSecret = map[string]string{c.tsigKey: c.tsigSec}

	reply, _, err := client.ExchangeContext(ctx, msg, c.server+":53")
	if err != nil {
		c.logger.Error("Failed to send DNS delete",
			zap.String("fqdn", fqdn),
			zap.String("server", c.server),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send DNS delete to %s: %w", c.server, err)
	}

	if reply.Rcode != dns.RcodeSuccess {
		c.logger.Error("DNS delete failed",
			zap.String("fqdn", fqdn),
			zap.String("server", c.server),
			zap.Int("rcode", reply.Rcode),
			zap.String("rcode_name", dns.RcodeToString[reply.Rcode]),
		)
		return fmt.Errorf("DNS delete failed: %s (rcode: %d)", dns.RcodeToString[reply.Rcode], reply.Rcode)
	}

	c.logger.Info("TXT record deleted successfully",
		zap.String("fqdn", fqdn),
		zap.String("server", c.server),
	)
	return nil
}

