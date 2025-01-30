package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// PDFFormExtractor implements FormExtractor for PDF files
type PDFFormExtractor struct {
	file *os.File
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

			cleanName, label := cleanPDFFieldName(fieldName)
			formFields = append(formFields, FormField{
				Name:  cleanName,
				Type:  fieldType,
				Label: label,
			})
		}
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

// Helper function to download a file
func downloadFile(url string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "form-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}

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

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error saving file: %w", err)
	}

	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("error seeking file: %w", err)
	}

	return tempFile, nil
}

// Helper function to clean PDF field names
func cleanPDFFieldName(name string) (string, string) {
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return name, name
	}

	baseName := strings.Join(parts[1:], " ")
	baseName = strings.TrimSpace(baseName)

	suffix := ""
	if idx := strings.LastIndex(baseName, "_"); idx != -1 {
		possibleSuffix := baseName[idx+1:]
		if _, err := strconv.Atoi(possibleSuffix); err == nil {
			suffix = baseName[idx:]
			baseName = strings.TrimSpace(baseName[:idx])
		}
	}

	baseName = strings.TrimSpace(baseName)
	baseName = strings.Trim(baseName, "\"")
	baseName = strings.Trim(baseName, "}")

	cleanName := baseName
	if suffix != "" {
		cleanName = baseName + suffix
	}

	label := baseName
	if suffix != "" {
		suffixNum := strings.TrimPrefix(suffix, "_")
		label = fmt.Sprintf("%s (%s)", baseName, suffixNum)
	}

	return cleanName, label
}
