package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/josephmowjew/form-field-extractor/pkg/scrapper"
)

func main() {
	url := flag.String("url", "", "URL of the form to extract (PDF or HTML)")
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout for operations")
	maxAttempts := flag.Int("max-attempts", 3, "Maximum number of retry attempts")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Create a new scrapper instance with options
	s := scrapper.New(
		scrapper.WithTimeout(*timeout),
		scrapper.WithMaxAttempts(*maxAttempts),
	)

	// Extract fields from the URL
	fields, err := s.ExtractFields(*url)
	if err != nil {
		log.Fatalf("Failed to extract fields: %v", err)
	}

	// Convert to JSON and print
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
