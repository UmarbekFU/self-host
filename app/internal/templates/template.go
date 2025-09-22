package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Template struct {
	Name        string
	Subject     string
	HTML        string
	Text        string
	Description string
}

type TemplateManager struct {
	templates map[string]*Template
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*Template),
	}
}

func (tm *TemplateManager) LoadTemplates(templatesDir string) error {
	// Load MJML templates and convert to HTML
	templates := []string{"welcome", "newsletter", "announcement"}
	
	for _, templateName := range templates {
		template, err := tm.loadTemplate(templatesDir, templateName)
		if err != nil {
			return fmt.Errorf("failed to load template %s: %w", templateName, err)
		}
		
		tm.templates[templateName] = template
	}
	
	return nil
}

func (tm *TemplateManager) loadTemplate(templatesDir, name string) (*Template, error) {
	// Load MJML file
	mjmlPath := filepath.Join(templatesDir, name+".mjml")
	mjmlContent, err := ioutil.ReadFile(mjmlPath)
	if err != nil {
		return nil, err
	}
	
	// For now, we'll use the MJML content as HTML
	// In production, you'd want to use a proper MJML compiler
	htmlContent := string(mjmlContent)
	
	// Generate text version (simplified)
	textContent := tm.generateTextVersion(htmlContent)
	
	// Generate subject from template
	subject := tm.generateSubject(name)
	
	return &Template{
		Name:        name,
		Subject:     subject,
		HTML:        htmlContent,
		Text:        textContent,
		Description: tm.getDescription(name),
	}, nil
}

func (tm *TemplateManager) generateTextVersion(html string) string {
	// Simple HTML to text conversion
	// In production, you'd want to use a proper HTML to text converter
	
	// Remove HTML tags
	text := html
	text = strings.ReplaceAll(text, "<mj-text>", "")
	text = strings.ReplaceAll(text, "</mj-text>", "\n")
	text = strings.ReplaceAll(text, "<mj-button>", "")
	text = strings.ReplaceAll(text, "</mj-button>", "")
	text = strings.ReplaceAll(text, "<mj-section>", "")
	text = strings.ReplaceAll(text, "</mj-section>", "\n")
	text = strings.ReplaceAll(text, "<mj-column>", "")
	text = strings.ReplaceAll(text, "</mj-column>", "")
	text = strings.ReplaceAll(text, "<mj-head>", "")
	text = strings.ReplaceAll(text, "</mj-head>", "")
	text = strings.ReplaceAll(text, "<mj-body>", "")
	text = strings.ReplaceAll(text, "</mj-body>", "")
	text = strings.ReplaceAll(text, "<mj-title>", "")
	text = strings.ReplaceAll(text, "</mj-title>", "")
	text = strings.ReplaceAll(text, "<mj-preview>", "")
	text = strings.ReplaceAll(text, "</mj-preview>", "")
	text = strings.ReplaceAll(text, "<mj-attributes>", "")
	text = strings.ReplaceAll(text, "</mj-attributes>", "")
	text = strings.ReplaceAll(text, "<mj-style>", "")
	text = strings.ReplaceAll(text, "</mj-style>", "")
	text = strings.ReplaceAll(text, "<mj-all>", "")
	text = strings.ReplaceAll(text, "</mj-all>", "")
	
	// Clean up extra whitespace
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}
	
	return strings.Join(cleanedLines, "\n")
}

func (tm *TemplateManager) generateSubject(templateName string) string {
	subjects := map[string]string{
		"welcome":      "Welcome to {{company_name}}!",
		"newsletter":   "{{newsletter_title}} - {{newsletter_date}}",
		"announcement": "Important: {{announcement_title}}",
	}
	
	if subject, exists := subjects[templateName]; exists {
		return subject
	}
	
	return "Newsletter from {{company_name}}"
}

func (tm *TemplateManager) getDescription(templateName string) string {
	descriptions := map[string]string{
		"welcome":      "Welcome email for new subscribers",
		"newsletter":   "Regular newsletter template with articles and updates",
		"announcement": "Important announcements and updates",
	}
	
	if desc, exists := descriptions[templateName]; exists {
		return desc
	}
	
	return "Email template"
}

func (tm *TemplateManager) GetTemplate(name string) (*Template, error) {
	template, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}
	
	return template, nil
}

func (tm *TemplateManager) ListTemplates() []*Template {
	var templates []*Template
	for _, template := range tm.templates {
		templates = append(templates, template)
	}
	return templates
}

func (tm *TemplateManager) RenderTemplate(templateName string, data map[string]interface{}) (*Template, error) {
	template, err := tm.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}
	
	// Create a copy of the template
	rendered := &Template{
		Name:        template.Name,
		Subject:     template.Subject,
		HTML:        template.HTML,
		Text:        template.Text,
		Description: template.Description,
	}
	
	// Render subject
	subjectTmpl, err := template.New("subject").Parse(template.Subject)
	if err == nil {
		var buf bytes.Buffer
		subjectTmpl.Execute(&buf, data)
		rendered.Subject = buf.String()
	}
	
	// Render HTML
	htmlTmpl, err := template.New("html").Parse(template.HTML)
	if err == nil {
		var buf bytes.Buffer
		htmlTmpl.Execute(&buf, data)
		rendered.HTML = buf.String()
	}
	
	// Render text
	textTmpl, err := template.New("text").Parse(template.Text)
	if err == nil {
		var buf bytes.Buffer
		textTmpl.Execute(&buf, data)
		rendered.Text = buf.String()
	}
	
	return rendered, nil
}

// Default template data
func GetDefaultTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"company_name":    "Your Company",
		"company_tagline": "Your tagline here",
		"company_address": "123 Main St, City, State 12345",
		"company_phone":   "+1 (555) 123-4567",
		"support_email":   "support@example.com",
		"support_phone":   "+1 (555) 123-4567",
		"first_name":      "{{first_name}}",
		"unsubscribe_url": "{{unsubscribe_url}}",
		"preferences_url": "{{preferences_url}}",
		"welcome_url":     "{{welcome_url}}",
	}
}
