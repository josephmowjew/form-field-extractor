package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// FormField represents a field in either a PDF or HTML form
type FormField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required,omitempty"`
	Value    string `json:"value,omitempty"`
}

// FormExtractor interface for different form extraction implementations
type FormExtractor interface {
	Extract() ([]FormField, error)
	Close() error
}

// PDFFormExtractor implements FormExtractor for PDF files
type PDFFormExtractor struct {
	file *os.File
}

// HTMLFormExtractor implements FormExtractor for HTML forms
type HTMLFormExtractor struct {
	browser *rod.Browser
	page    *rod.Page
	url     string
	timeout time.Duration
}

// Config holds the application configuration
type Config struct {
	URL         string
	Timeout     time.Duration
	OutputPath  string
	MaxAttempts int
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		URL:         "https://www.txdmv.gov/sites/default/files/form_files/130-U.pdf",
		Timeout:     30 * time.Second,
		OutputPath:  "output.pdf",
		MaxAttempts: 3,
	}
}

// NewFormExtractor creates the appropriate form extractor based on the URL
func NewFormExtractor(url string) (FormExtractor, error) {
	if strings.HasSuffix(strings.ToLower(url), ".pdf") {
		return NewPDFFormExtractor(url)
	}
	return NewHTMLFormExtractor(url)
}

// NewPDFFormExtractor creates a new PDF form extractor
func NewPDFFormExtractor(url string) (*PDFFormExtractor, error) {
	tempFile, err := downloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}

	return &PDFFormExtractor{
		file: tempFile,
	}, nil
}

// NewHTMLFormExtractor creates a new HTML form extractor
func NewHTMLFormExtractor(url string) (*HTMLFormExtractor, error) {
	// Configure browser with timeout
	launcher := rod.New()
	browser := launcher.MustConnect()

	// Create new page
	page := browser.MustPage()

	// Navigate with timeout and error handling
	err := page.Timeout(10 * time.Second).Navigate(url)
	if err != nil {
		browser.MustClose()
		return nil, fmt.Errorf("failed to navigate to URL: %w", err)
	}

	// Wait for page load with timeout
	err = page.Timeout(5 * time.Second).WaitLoad()
	if err != nil {
		browser.MustClose()
		return nil, fmt.Errorf("timeout waiting for page to load: %w", err)
	}

	return &HTMLFormExtractor{
		browser: browser,
		page:    page,
		url:     url,
		timeout: 30 * time.Second,
	}, nil
}

// Extract implements FormExtractor for PDFFormExtractor
func (p *PDFFormExtractor) Extract() ([]FormField, error) {
	fields, err := api.FormFields(p.file, nil)
	if err != nil {
		return nil, fmt.Errorf("error listing form fields: %w", err)
	}

	var formFields []FormField

	for _, field := range fields {
		parts := strings.Fields(fmt.Sprintf("%v", field))
		if len(parts) >= 3 {
			fieldType := parts[2]
			fieldName := strings.Join(parts[3:], " ")

			// Clean up the field name
			cleanName, label := cleanPDFFieldName(fieldName)

			// For PDF forms, we'll use the cleaned name as both name and label
			formFields = append(formFields, FormField{
				Name:  cleanName,
				Type:  fieldType,
				Label: label,
			})
		}
	}

	return formFields, nil
}

// Extract implements FormExtractor for HTMLFormExtractor
func (h *HTMLFormExtractor) Extract() ([]FormField, error) {
	var formFields []FormField

	// Try to find form elements with timeout
	elements, err := h.page.Timeout(5 * time.Second).Elements("input, select, textarea")
	if err != nil {
		return nil, fmt.Errorf("failed to find form elements: %w", err)
	}

	for _, element := range elements {
		// Get element attributes with error handling
		typeStr := "text" // default type
		if t, err := element.Attribute("type"); err == nil && t != nil {
			typeStr = *t
		}

		// Get name attribute
		name, err := element.Attribute("name")
		if err != nil || name == nil {
			continue // Skip elements without name
		}

		// Try to find label
		label := ""

		// First try to find associated label using 'for' attribute
		if id, err := element.Attribute("id"); err == nil && id != nil {
			if labelElement, err := h.page.Element(fmt.Sprintf("label[for='%s']", *id)); err == nil && labelElement != nil {
				if text, err := labelElement.Text(); err == nil {
					label = text
				}
			}
		}

		// If no label found, try aria-label
		if label == "" {
			if ariaLabel, err := element.Attribute("aria-label"); err == nil && ariaLabel != nil {
				label = *ariaLabel
			}
		}

		// If still no label, try placeholder
		if label == "" {
			if placeholder, err := element.Attribute("placeholder"); err == nil && placeholder != nil {
				label = *placeholder
			}
		}

		// If still no label, use name
		if label == "" {
			label = *name
		}

		// Check if field is required
		required := false
		if req, err := element.Attribute("required"); err == nil && req != nil {
			required = true
		}

		formFields = append(formFields, FormField{
			Name:     *name,
			Type:     typeStr,
			Label:    label,
			Required: required,
		})
	}

	return formFields, nil
}

// Close implements FormExtractor for PDFFormExtractor
func (p *PDFFormExtractor) Close() error {
	if p.file != nil {
		if err := p.file.Close(); err != nil {
			return fmt.Errorf("error closing PDF file: %w", err)
		}
		if err := os.Remove(p.file.Name()); err != nil {
			return fmt.Errorf("error removing temporary PDF file: %w", err)
		}
	}
	return nil
}

// Close implements FormExtractor for HTMLFormExtractor
func (h *HTMLFormExtractor) Close() error {
	if h.browser != nil {
		h.browser.MustClose()
	}
	return nil
}

func main() {
	if err := run(DefaultConfig()); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run(cfg *Config) error {
	extractor, err := NewFormExtractor(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to create form extractor: %w", err)
	}
	defer extractor.Close()

	fields, err := extractor.Extract()
	if err != nil {
		return fmt.Errorf("failed to extract fields: %w", err)
	}

	// Convert to JSON and print
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

func downloadFile(url string) (*os.File, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "form-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copy the body to the file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error saving file: %w", err)
	}

	// Seek to beginning of file for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error seeking file: %w", err)
	}

	return tempFile, nil
}

// Helper function to clean PDF field names
func cleanPDFFieldName(name string) (string, string) {
	// Split into parts
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return name, name // Return as is if can't split
	}

	// The first part is typically the numeric ID, remove it
	// Join the rest to get the base name
	baseName := strings.Join(parts[1:], " ")
	baseName = strings.TrimSpace(baseName)

	// Extract any numeric suffix and clean the base name
	suffix := ""
	if idx := strings.LastIndex(baseName, "_"); idx != -1 {
		possibleSuffix := baseName[idx+1:]
		if _, err := strconv.Atoi(possibleSuffix); err == nil {
			// If the part after _ is a number, it's a suffix
			suffix = baseName[idx:]
			baseName = strings.TrimSpace(baseName[:idx])
		}
	}

	// Clean up any remaining special characters or extra spaces
	baseName = strings.TrimSpace(baseName)
	baseName = strings.Trim(baseName, "\"")
	baseName = strings.Trim(baseName, "}")

	// Create the clean name (preserving suffix for uniqueness)
	cleanName := baseName
	if suffix != "" {
		cleanName = baseName + suffix
	}

	// Create a human-readable label
	label := baseName
	if suffix != "" {
		// Convert _2 to (2) for better readability
		suffixNum := strings.TrimPrefix(suffix, "_")
		label = fmt.Sprintf("%s (%s)", baseName, suffixNum)
	}

	return cleanName, label
}
