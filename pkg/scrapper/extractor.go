package scrapper

import (
	"fmt"
	"strings"
)

// FormExtractor interface for different form extraction implementations
type FormExtractor interface {
	Extract() ([]FormField, error)
	Close() error
}

// Scrapper is the main library interface
type Scrapper struct {
	config *Config
}

// New creates a new Scrapper instance with the provided options
func New(options ...Option) *Scrapper {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return &Scrapper{
		config: config,
	}
}

// ExtractFields extracts form fields from the provided URL
func (s *Scrapper) ExtractFields(url string) ([]FormField, error) {
	extractor, err := s.newFormExtractor(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create form extractor: %w", err)
	}
	defer extractor.Close()

	fields, err := extractor.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract fields: %w", err)
	}

	return fields, nil
}

// newFormExtractor creates the appropriate form extractor based on the URL
func (s *Scrapper) newFormExtractor(url string) (FormExtractor, error) {
	if strings.HasSuffix(strings.ToLower(url), ".pdf") {
		return NewPDFFormExtractor(url)
	}
	return NewHTMLFormExtractor(url, s.config.Timeout)
}
