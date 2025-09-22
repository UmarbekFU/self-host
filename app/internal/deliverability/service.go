package deliverability

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"newsletter/internal/store"
	"github.com/sirupsen/logrus"
)

type Service struct {
	// Add any configuration or dependencies here
}

func NewService() *Service {
	return &Service{}
}

type DomainStatus struct {
	Domain     string            `json:"domain"`
	SPF        CheckResult       `json:"spf"`
	DKIM       CheckResult       `json:"dkim"`
	DMARC      CheckResult       `json:"dmarc"`
	PTR        CheckResult       `json:"ptr"`
	TLS        CheckResult       `json:"tls"`
	Overall    string            `json:"overall"`
	Checks     map[string]bool   `json:"checks"`
}

type CheckResult struct {
	Status    string `json:"status"`    // "pass", "fail", "warning"
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

func (s *Service) CheckDomainStatus(domain *store.Domain) (*DomainStatus, error) {
	status := &DomainStatus{
		Domain: domain.Domain,
		Checks: make(map[string]bool),
	}

	// Check SPF record
	spfResult, err := s.checkSPF(domain.Domain, domain.SPFRecord)
	if err != nil {
		logrus.Errorf("SPF check failed for %s: %v", domain.Domain, err)
		spfResult = CheckResult{Status: "fail", Message: "SPF check failed", Details: err.Error()}
	}
	status.SPF = spfResult
	status.Checks["spf"] = spfResult.Status == "pass"

	// Check DKIM record
	dkimResult, err := s.checkDKIM(domain.Domain, domain.DKIMSelector, domain.DKIMPublicKey)
	if err != nil {
		logrus.Errorf("DKIM check failed for %s: %v", domain.Domain, err)
		dkimResult = CheckResult{Status: "fail", Message: "DKIM check failed", Details: err.Error()}
	}
	status.DKIM = dkimResult
	status.Checks["dkim"] = dkimResult.Status == "pass"

	// Check DMARC record
	dmarcResult, err := s.checkDMARC(domain.Domain, domain.DMARCRecord)
	if err != nil {
		logrus.Errorf("DMARC check failed for %s: %v", domain.Domain, err)
		dmarcResult = CheckResult{Status: "fail", Message: "DMARC check failed", Details: err.Error()}
	}
	status.DMARC = dmarcResult
	status.Checks["dmarc"] = dmarcResult.Status == "pass"

	// Check PTR record
	ptrResult, err := s.checkPTR(domain.PTRRecord)
	if err != nil {
		logrus.Errorf("PTR check failed for %s: %v", domain.Domain, err)
		ptrResult = CheckResult{Status: "fail", Message: "PTR check failed", Details: err.Error()}
	}
	status.PTR = ptrResult
	status.Checks["ptr"] = ptrResult.Status == "pass"

	// Check TLS (simplified)
	tlsResult := s.checkTLS(domain.Domain)
	status.TLS = tlsResult
	status.Checks["tls"] = tlsResult.Status == "pass"

	// Determine overall status
	status.Overall = s.determineOverallStatus(status.Checks)

	return status, nil
}

func (s *Service) checkSPF(domain, expectedSPF string) (CheckResult, error) {
	// Look up SPF record
	records, err := net.LookupTXT(domain)
	if err != nil {
		return CheckResult{Status: "fail", Message: "Failed to lookup TXT records"}, err
	}

	var spfRecord string
	for _, record := range records {
		if strings.HasPrefix(record, "v=spf1") {
			spfRecord = record
			break
		}
	}

	if spfRecord == "" {
		return CheckResult{Status: "fail", Message: "No SPF record found"}, nil
	}

	// Compare with expected record
	if spfRecord != expectedSPF {
		return CheckResult{
			Status:  "warning",
			Message: "SPF record doesn't match expected value",
			Details: fmt.Sprintf("Expected: %s\nFound: %s", expectedSPF, spfRecord),
		}, nil
	}

	return CheckResult{Status: "pass", Message: "SPF record is valid"}, nil
}

func (s *Service) checkDKIM(domain, selector, publicKey string) (CheckResult, error) {
	// Look up DKIM record
	dkimDomain := fmt.Sprintf("%s._domainkey.%s", selector, domain)
	records, err := net.LookupTXT(dkimDomain)
	if err != nil {
		return CheckResult{Status: "fail", Message: "Failed to lookup DKIM record"}, err
	}

	if len(records) == 0 {
		return CheckResult{Status: "fail", Message: "No DKIM record found"}, nil
	}

	// Check if the public key matches (simplified check)
	dkimRecord := records[0]
	if !strings.Contains(dkimRecord, "v=DKIM1") {
		return CheckResult{Status: "fail", Message: "Invalid DKIM record format"}, nil
	}

	// In a real implementation, you'd parse the DKIM record and compare the public key
	// For now, we'll just check if the record exists and has the right format
	return CheckResult{Status: "pass", Message: "DKIM record is valid"}, nil
}

func (s *Service) checkDMARC(domain, expectedDMARC string) (CheckResult, error) {
	// Look up DMARC record
	dmarcDomain := fmt.Sprintf("_dmarc.%s", domain)
	records, err := net.LookupTXT(dmarcDomain)
	if err != nil {
		return CheckResult{Status: "fail", Message: "Failed to lookup DMARC record"}, err
	}

	if len(records) == 0 {
		return CheckResult{Status: "fail", Message: "No DMARC record found"}, nil
	}

	dmarcRecord := records[0]
	if !strings.HasPrefix(dmarcRecord, "v=DMARC1") {
		return CheckResult{Status: "fail", Message: "Invalid DMARC record format"}, nil
	}

	// Compare with expected record
	if dmarcRecord != expectedDMARC {
		return CheckResult{
			Status:  "warning",
			Message: "DMARC record doesn't match expected value",
			Details: fmt.Sprintf("Expected: %s\nFound: %s", expectedDMARC, dmarcRecord),
		}, nil
	}

	return CheckResult{Status: "pass", Message: "DMARC record is valid"}, nil
}

func (s *Service) checkPTR(expectedPTR string) (CheckResult, error) {
	// Get server's public IP
	serverIP, err := s.getServerIP()
	if err != nil {
		return CheckResult{Status: "fail", Message: "Failed to get server IP"}, err
	}

	// Perform reverse DNS lookup
	names, err := net.LookupAddr(serverIP)
	if err != nil {
		return CheckResult{Status: "fail", Message: "Failed to perform reverse DNS lookup"}, err
	}

	if len(names) == 0 {
		return CheckResult{Status: "fail", Message: "No PTR record found"}, nil
	}

	// Check if any of the returned names match the expected PTR record
	for _, name := range names {
		// Remove trailing dot
		name = strings.TrimSuffix(name, ".")
		if name == expectedPTR {
			return CheckResult{Status: "pass", Message: "PTR record is valid"}, nil
		}
	}

	return CheckResult{
		Status:  "warning",
		Message: "PTR record doesn't match expected value",
		Details: fmt.Sprintf("Expected: %s\nFound: %s", expectedPTR, names[0]),
	}, nil
}

func (s *Service) getServerIP() (string, error) {
	// Try to get external IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func (s *Service) checkTLS(domain string) CheckResult {
	// Check if domain has a valid TLS certificate
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return CheckResult{Status: "fail", Message: "TLS connection failed", Details: err.Error()}
	}
	defer conn.Close()

	// Check certificate validity
	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return CheckResult{Status: "fail", Message: "No TLS certificate found"}
	}

	cert := state.PeerCertificates[0]
	now := time.Now()
	
	if now.Before(cert.NotBefore) {
		return CheckResult{Status: "fail", Message: "TLS certificate not yet valid"}
	}
	
	if now.After(cert.NotAfter) {
		return CheckResult{Status: "fail", Message: "TLS certificate expired"}
	}

	// Check if certificate matches domain
	if err := cert.VerifyHostname(domain); err != nil {
		return CheckResult{Status: "warning", Message: "TLS certificate doesn't match domain", Details: err.Error()}
	}

	return CheckResult{Status: "pass", Message: "TLS certificate is valid"}
}

func (s *Service) determineOverallStatus(checks map[string]bool) string {
	// Count passed checks
	total := len(checks)
	passedCount := 0
	
	for _, passed := range checks {
		if passed {
			passedCount++
		}
	}
	
	// Determine overall status based on percentage of passed checks
	percentage := float64(passedCount) / float64(total) * 100
	
	if percentage >= 80 {
		return "pass"
	} else if percentage >= 60 {
		return "warning"
	} else {
		return "fail"
	}
}

func (s *Service) ValidateEmailContent(html, text string) []string {
	var warnings []string
	
	// Check if text version exists
	if strings.TrimSpace(text) == "" {
		warnings = append(warnings, "No text version provided")
	}
	
	// Check HTML to text ratio
	htmlLength := len(strings.TrimSpace(html))
	textLength := len(strings.TrimSpace(text))
	
	if htmlLength > 0 && textLength > 0 {
		ratio := float64(htmlLength) / float64(textLength)
		if ratio > 3.0 {
			warnings = append(warnings, "HTML to text ratio is too high (may trigger spam filters)")
		}
	}
	
	// Check for spammy words (simplified)
	spamWords := []string{"free", "win", "winner", "congratulations", "urgent", "act now", "limited time"}
	content := strings.ToLower(html + " " + text)
	
	for _, word := range spamWords {
		if strings.Contains(content, word) {
			warnings = append(warnings, fmt.Sprintf("Content contains potentially spammy word: %s", word))
		}
	}
	
	// Check for excessive links
	linkCount := strings.Count(html, "<a href")
	if linkCount > 10 {
		warnings = append(warnings, "Too many links in email (may trigger spam filters)")
	}
	
	return warnings
}

func (s *Service) GetRequiredDNSRecords(domain *store.Domain) map[string]string {
	records := make(map[string]string)
	
	// SPF record
	records["SPF"] = domain.SPFRecord
	
	// DKIM record
	dkimDomain := fmt.Sprintf("%s._domainkey.%s", domain.DKIMSelector, domain.Domain)
	records["DKIM"] = fmt.Sprintf("%s: %s", dkimDomain, domain.DKIMPublicKey)
	
	// DMARC record
	dmarcDomain := fmt.Sprintf("_dmarc.%s", domain.Domain)
	records["DMARC"] = fmt.Sprintf("%s: %s", dmarcDomain, domain.DMARCRecord)
	
	// PTR record
	records["PTR"] = fmt.Sprintf("Reverse DNS should point to: %s", domain.PTRRecord)
	
	return records
}
