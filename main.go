package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Config holds the application configuration
type Config struct {
	PDFURL      string
	Timeout     time.Duration
	OutputPath  string
	MaxAttempts int
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		PDFURL:      "https://www.txdmv.gov/sites/default/files/form_files/130-U.pdf",
		Timeout:     30 * time.Second,
		OutputPath:  "output.pdf",
		MaxAttempts: 3,
	}
}

func main() {
	if err := run(DefaultConfig()); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run(cfg *Config) error {
	// Download PDF to temporary file
	tempPath, err := downloadPDF(cfg.PDFURL)
	if err != nil {
		return fmt.Errorf("failed to download PDF: %w", err)
	}
	defer cleanupFile(tempPath)

	// Extract and process form fields
	if err := extractPDFFields(tempPath); err != nil {
		return fmt.Errorf("failed to extract PDF fields: %w", err)
	}

	return nil
}

func downloadPDF(url string) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "form-*.pdf")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer tempFile.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading PDF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copy the body to the file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving PDF: %w", err)
	}

	return tempFile.Name(), nil
}

func cleanupFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Warning: failed to remove temporary file %s: %v", path, err)
	}
}

func extractPDFFields(pdfPath string) error {
	// Open the PDF file
	f, err := os.Open(pdfPath)
	if err != nil {
		return fmt.Errorf("error opening PDF: %w", err)
	}
	defer f.Close()

	// Get form fields
	fields, err := api.FormFields(f, nil)
	if err != nil {
		return fmt.Errorf("error listing form fields: %w", err)
	}

	type FieldInfo struct {
		Name string `json:"field_name"`
		Type string `json:"field_type"`
	}

	type FormFields struct {
		Fields []FieldInfo `json:"fields"`
	}

	var formFields FormFields

	// Process form fields
	for _, field := range fields {
		parts := strings.Fields(fmt.Sprintf("%v", field))
		if len(parts) >= 3 {
			fieldType := parts[2]
			fieldName := strings.Join(parts[3:], " ")
			cleanName := fieldName
			if idx := strings.Index(fieldName, " "); idx != -1 {
				cleanName = strings.TrimSpace(fieldName[idx:])
			}
			cleanName = strings.TrimSuffix(cleanName, " }")

			formFields.Fields = append(formFields.Fields, FieldInfo{
				Name: cleanName,
				Type: fieldType,
			})
		}
	}

	// Convert to JSON and print
	jsonData, err := json.MarshalIndent(formFields, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
