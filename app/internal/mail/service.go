package mail

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"newsletter/internal/store"
)

type Service struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func NewService() *Service {
	return &Service{
		SMTPHost:     "localhost",
		SMTPPort:     "587",
		SMTPUsername: "",
		SMTPPassword: "",
	}
}

type Message struct {
	To          []string
	From        string
	FromName    string
	Subject     string
	HTML        string
	Text        string
	ReplyTo     string
	Headers     map[string]string
	DKIMDomain  string
	DKIMKey     string
	DKIMSelector string
}

func (s *Service) Send(msg *Message) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", msg.FromName, msg.From)
	e.To = msg.To
	e.Subject = msg.Subject
	e.HTML = []byte(msg.HTML)
	e.Text = []byte(msg.Text)

	if msg.ReplyTo != "" {
		e.ReplyTo = []string{msg.ReplyTo}
	}

	// Add custom headers
	for key, value := range msg.Headers {
		e.Headers.Set(key, value)
	}

	// Add DKIM signature if provided
	if msg.DKIMDomain != "" && msg.DKIMKey != "" && msg.DKIMSelector != "" {
		if err := s.addDKIMSignature(e, msg); err != nil {
			logrus.Warnf("Failed to add DKIM signature: %v", err)
		}
	}

	// Add List-Unsubscribe header
	unsubscribeURL := fmt.Sprintf("https://%s/u/%%recipient_id%%/%%unsubscribe_token%%", msg.DKIMDomain)
	e.Headers.Set("List-Unsubscribe", fmt.Sprintf("<%s>", unsubscribeURL))
	e.Headers.Set("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")

	// Send email
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)
	
	return e.Send(addr, auth)
}

func (s *Service) addDKIMSignature(e *email.Email, msg *Message) error {
	// Parse private key
	block, _ := pem.Decode([]byte(msg.DKIMKey))
	if block == nil {
		return fmt.Errorf("failed to decode DKIM private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse DKIM private key: %w", err)
	}

	// Create DKIM signature
	signature := s.createDKIMSignature(e, msg.DKIMDomain, msg.DKIMSelector, privateKey)
	
	// Add DKIM-Signature header
	e.Headers.Set("DKIM-Signature", signature)
	
	return nil
}

func (s *Service) createDKIMSignature(e *email.Email, domain, selector string, privateKey *rsa.PrivateKey) string {
	// This is a simplified DKIM implementation
	// In production, you'd want to use a proper DKIM library
	
	// Convert MIMEHeader to map[string]string
	headers := make(map[string]string)
	for k, v := range e.Headers {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	
	canonicalizedHeaders := s.canonicalizeHeaders(headers)
	canonicalizedBody := s.canonicalizeBody(e.Text)
	
	// Create the signature data (using canonicalized data in real implementation)
	signatureData := fmt.Sprintf("v=1; a=rsa-sha256; c=relaxed/relaxed; d=%s; s=%s; h=%s; bh=%s; b=",
		domain, selector, "from:to:subject:date:message-id", "dummy-hash")
	
	// In a real implementation, you'd sign the canonicalized data
	// For now, we'll return a placeholder signature
	_ = canonicalizedHeaders // Use variable to avoid "declared and not used" error
	_ = canonicalizedBody     // Use variable to avoid "declared and not used" error
	return signatureData + "dummy-signature"
}

func (s *Service) canonicalizeHeaders(headers map[string]string) string {
	// Simplified header canonicalization
	var result []string
	for key, value := range headers {
		result = append(result, fmt.Sprintf("%s: %s", strings.ToLower(key), strings.TrimSpace(value)))
	}
	return strings.Join(result, "\r\n")
}

func (s *Service) canonicalizeBody(body []byte) string {
	// Simplified body canonicalization
	bodyStr := string(body)
	// Remove trailing whitespace from each line
	lines := strings.Split(bodyStr, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func GenerateDKIMKeys() (privateKey, publicKey string, err error) {
	// Generate RSA key pair
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Encode private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	// Encode public key
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

func (s *Service) CreateTestMessage(campaignID int, subscriberID int, to string) *Message {
	return &Message{
		To:       []string{to},
		From:     "test@example.com",
		FromName: "Test Sender",
		Subject:  "Test Campaign",
		HTML:     "<p>This is a test email.</p>",
		Text:     "This is a test email.",
		Headers: map[string]string{
			"X-Campaign-ID":      fmt.Sprintf("%d", campaignID),
			"X-Subscriber-ID":    fmt.Sprintf("%d", subscriberID),
			"X-Mailer":           "Newsletter Platform",
		},
	}
}

func (s *Service) CreateCampaignMessage(campaign *store.Campaign, subscriber *store.Subscriber) *Message {
	// Replace placeholders in HTML and text
	html := s.replacePlaceholders(campaign.HTML, subscriber)
	text := s.replacePlaceholders(campaign.Text, subscriber)

	// Add tracking pixels and links
	html = s.addTracking(html, campaign.ID, subscriber.ID)

	return &Message{
		To:       []string{subscriber.Email},
		From:     campaign.FromEmail,
		FromName: campaign.FromName,
		Subject:  campaign.Subject,
		HTML:     html,
		Text:     text,
		ReplyTo:  campaign.ReplyTo,
		Headers: map[string]string{
			"X-Campaign-ID":      fmt.Sprintf("%d", campaign.ID),
			"X-Subscriber-ID":    fmt.Sprintf("%d", subscriber.ID),
			"X-Mailer":           "Newsletter Platform",
			"Message-ID":         fmt.Sprintf("<%d.%d@newsletter.local>", campaign.ID, subscriber.ID),
		},
	}
}

func (s *Service) replacePlaceholders(content string, subscriber *store.Subscriber) string {
	// Replace common placeholders
	content = strings.ReplaceAll(content, "{{email}}", subscriber.Email)
	content = strings.ReplaceAll(content, "{{unsubscribe_url}}", fmt.Sprintf("https://example.com/u/%d/token", subscriber.ID))
	
	// Replace custom attributes
	if subscriber.Attributes != nil {
		var attrs map[string]interface{}
		if err := json.Unmarshal(subscriber.Attributes, &attrs); err == nil {
			for key, value := range attrs {
				placeholder := fmt.Sprintf("{{%s}}", key)
				content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
			}
		}
	}
	
	return content
}

func (s *Service) addTracking(html string, campaignID, subscriberID int) string {
	// Add open tracking pixel
	trackingPixel := fmt.Sprintf(`<img src="https://example.com/api/track/open?c=%d&s=%d" width="1" height="1" style="display:none;">`, campaignID, subscriberID)
	
	// Add before closing body tag
	if strings.Contains(html, "</body>") {
		html = strings.Replace(html, "</body>", trackingPixel+"</body>", 1)
	} else {
		html += trackingPixel
	}
	
	// TODO: Add click tracking for links
	// This would involve wrapping all links with tracking URLs
	
	return html
}
