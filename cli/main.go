package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sbro101/goresolve/v2"
)

// main is the entry point for the goresolve command line interface.
//
// Usage: goresolve <url> [nameserver]
//
// Example: goresolve example.com 8.8.8.8
//
// The goresolve command line interface looks up the IP addresses for the given
// hostname and prints them as a JSON object.
//
// The nameserver parameter is optional and defaults to 1.1.1.1 if not
// provided.
func main() {
	// Check the command line arguments.
	if len(os.Args) < 2 {
		fmt.Println("Usage: goresolve <url> [nameserver]")
		fmt.Println("Example: goresolve example.com 8.8.8.8")
		os.Exit(1)
	}

	// Set the hostname and nameserver from the command line arguments.
	url := os.Args[1]
	nameserver := "1.1.1.1" // Default nameserver

	// If the nameserver argument is provided, use it instead of the default.
	if len(os.Args) > 2 {
		nameserver = os.Args[2]
	}

	// Look up the hostname.
	rd, err := goresolve.Hostname(url, nameserver)
	if err != nil {
		// If there is an error, print an error message and exit.
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if rd == nil {
		fmt.Println("Error: nil pointer reference")
		os.Exit(1)
	}

	// Create a struct to hold the output data.
	output := struct {
		URL        string          `json:"url"`
		Nameserver string          `json:"nameserver"`
		Result     *goresolve.Data `json:"result"`
	}{
		URL:        url,
		Nameserver: nameserver,
		Result:     rd,
	}

	// Marshal the output struct to JSON.
	json, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		// If there is an error marshaling JSON, print an error message and exit.
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Print the output JSON.
	fmt.Printf("%s\n", json)
}
