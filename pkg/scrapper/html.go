package scrapper

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// HTMLFormExtractor implements FormExtractor for HTML forms
type HTMLFormExtractor struct {
	browser *rod.Browser
	page    *rod.Page
	url     string
	timeout time.Duration
}

// NewHTMLFormExtractor creates a new HTML form extractor
func NewHTMLFormExtractor(url string, timeout time.Duration) (*HTMLFormExtractor, error) {
	launcher := rod.New()
	browser := launcher.MustConnect()

	page := browser.MustPage()

	err := page.Timeout(timeout).Navigate(url)
	if err != nil {
		browser.MustClose()
		return nil, fmt.Errorf("failed to navigate to URL: %w", err)
	}

	err = page.Timeout(timeout).WaitLoad()
	if err != nil {
		browser.MustClose()
		return nil, fmt.Errorf("timeout waiting for page to load: %w", err)
	}

	return &HTMLFormExtractor{
		browser: browser,
		page:    page,
		url:     url,
		timeout: timeout,
	}, nil
}

// Extract implements FormExtractor for HTMLFormExtractor
func (h *HTMLFormExtractor) Extract() ([]FormField, error) {
	var formFields []FormField

	elements, err := h.page.Timeout(h.timeout).Elements("input, select, textarea")
	if err != nil {
		return nil, fmt.Errorf("failed to find form elements: %w", err)
	}

	for _, element := range elements {
		typeStr := "text" // default type
		if t, err := element.Attribute("type"); err == nil && t != nil {
			typeStr = *t
		}

		name, err := element.Attribute("name")
		if err != nil || name == nil {
			continue
		}

		label := ""

		if id, err := element.Attribute("id"); err == nil && id != nil {
			if labelElement, err := h.page.Element(fmt.Sprintf("label[for='%s']", *id)); err == nil && labelElement != nil {
				if text, err := labelElement.Text(); err == nil {
					label = text
				}
			}
		}

		if label == "" {
			if ariaLabel, err := element.Attribute("aria-label"); err == nil && ariaLabel != nil {
				label = *ariaLabel
			}
		}

		if label == "" {
			if placeholder, err := element.Attribute("placeholder"); err == nil && placeholder != nil {
				label = *placeholder
			}
		}

		if label == "" {
			label = *name
		}

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

// Close implements FormExtractor for HTMLFormExtractor
func (h *HTMLFormExtractor) Close() error {
	if h.browser != nil {
		h.browser.MustClose()
	}
	return nil
}
