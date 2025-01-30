package scrapper

import "time"

// FormField represents a field in either a PDF or HTML form
type FormField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Label    string `json:"label"`
	Required bool   `json:"required,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Config holds the configuration for the form extractor
type Config struct {
	Timeout     time.Duration
	MaxAttempts int
}

// Option is a function that modifies the Config
type Option func(*Config)

// WithTimeout sets the timeout for operations
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithMaxAttempts sets the maximum number of retry attempts
func WithMaxAttempts(attempts int) Option {
	return func(c *Config) {
		c.MaxAttempts = attempts
	}
}

// defaultConfig returns the default configuration
func defaultConfig() *Config {
	return &Config{
		Timeout:     30 * time.Second,
		MaxAttempts: 3,
	}
}
