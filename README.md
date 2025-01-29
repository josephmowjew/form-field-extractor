# Form Field Extractor

A powerful Go application that extracts form fields from both PDF files and HTML web forms. This tool automatically detects the input type (PDF or HTML) and returns a standardized JSON output of all form fields, including their names, types, labels, and required status.

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

```bash
# Clone the repository
git clone [your-repo-url]
cd [your-repo-name]

# Install dependencies
go mod download
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

func main() {
    cfg := DefaultConfig()
    
    // For PDF forms
    cfg.URL = "https://example.com/form.pdf"
    
    // Or for HTML forms
    // cfg.URL = "https://example.com/form.html"
    
    if err := run(cfg); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}
```

### Configuration Options

```go
type Config struct {
    URL         string        // URL of the PDF or HTML form
    Timeout     time.Duration // Timeout for operations
    OutputPath  string        // Path for output (if needed)
    MaxAttempts int          // Maximum retry attempts
}
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

### FormExtractor Interface

```go
type FormExtractor interface {
    Extract() ([]FormField, error)
    Close() error
}
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