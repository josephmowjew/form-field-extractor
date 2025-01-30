# Form Field Extractor

A powerful Go library and CLI tool that extracts form fields from both PDF files and HTML web forms. This tool automatically detects the input type (PDF or HTML) and returns a standardized JSON output of all form fields, including their names, types, labels, and required status.

[![Go Reference](https://pkg.go.dev/badge/github.com/josephmowjew/form-field-extractor.svg)](https://pkg.go.dev/github.com/josephmowjew/form-field-extractor)
[![Go Report Card](https://goreportcard.com/badge/github.com/josephmowjew/form-field-extractor)](https://goreportcard.com/report/github.com/josephmowjew/form-field-extractor)
![Version](https://img.shields.io/badge/version-v0.1.1-blue.svg)

## Features

- 🔄 Unified extraction interface for both PDF and HTML forms
- 📄 PDF form field extraction using pdfcpu
- 🌐 HTML form field extraction using Rod (headless browser)
- 🏷️ Intelligent label detection for HTML forms
- 🔍 Automatic field type detection
- 🧹 Smart field name cleaning and normalization
- ⚡ Concurrent processing capabilities
- ⏱️ Configurable timeouts and retry attempts
- 🎯 Support for various input types (text, select, textarea, etc.)

## Requirements

- Go 1.22 or newer
- Chrome/Chromium (for HTML form extraction)

## Installation

### As a Library

```bash
# Latest version
go get github.com/josephmowjew/form-field-extractor

# Specific version
go get github.com/josephmowjew/form-field-extractor@v0.1.0
```

### As a CLI Tool

```bash
# Clone the repository
git clone https://github.com/josephmowjew/form-field-extractor.git
cd form-field-extractor

# Checkout specific version (optional)
git checkout v0.1.0

# Install dependencies
go mod download

# Build the CLI tool
go build -o scrapper ./cmd/scrapper
```

## Usage

### Library Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/josephmowjew/form-field-extractor/pkg/scrapper"
)

func main() {
    // Create a new scrapper instance with options
    s := scrapper.New(
        scrapper.WithTimeout(30 * time.Second),
        scrapper.WithMaxAttempts(3),
    )
    
    // Extract fields from a URL (PDF or HTML)
    fields, err := s.ExtractFields("https://example.com/form.pdf")
    if err != nil {
        log.Fatalf("Failed to extract fields: %v", err)
    }
    
    // Convert to JSON
    jsonData, err := json.MarshalIndent(fields, "", "  ")
    if err != nil {
        log.Fatalf("Failed to marshal JSON: %v", err)
    }
    
    fmt.Println(string(jsonData))
}
```

### CLI Usage

```bash
# Extract fields from a PDF form
./scrapper -url https://example.com/form.pdf

# Extract fields from an HTML form with custom timeout
./scrapper -url https://example.com/form.html -timeout 45s

# Show help
./scrapper -help
```

### Configuration Options

The library supports the following configuration options:

```go
// Available options
scrapper.WithTimeout(30 * time.Second)  // Set operation timeout
scrapper.WithMaxAttempts(3)             // Set maximum retry attempts
```

### Example Output

```json
[
  {
    "name": "firstName",
    "type": "text",
    "label": "First Name",
    "required": true
  },
  {
    "name": "documentType",
    "type": "select",
    "label": "Document Type",
    "required": false
  }
]
```

## API Documentation

### Main Interface

```go
type Scrapper interface {
    ExtractFields(url string) ([]FormField, error)
}

// Create a new scrapper instance
scrapper := New(options ...Option)
```

### FormField Structure

```go
type FormField struct {
    Name     string `json:"name"`
    Type     string `json:"type"`
    Label    string `json:"label"`
    Required bool   `json:"required,omitempty"`
    Value    string `json:"value,omitempty"`
}
```

### Supported Form Field Types

- Text inputs
- Select dropdowns
- Textareas
- Checkboxes
- Radio buttons
- And more, depending on the form type (PDF/HTML)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) for PDF processing
- [Rod](https://github.com/go-rod/rod) for headless browser automation 