package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/schraf/assistant/internal/config"
	"github.com/schraf/assistant/internal/telegraph"
)

func main() {
	if err := config.LoadEnv(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Get default values from environment variables
	defaultShortName := os.Getenv("TELEGRAPH_SHORT_NAME")
	defaultAuthorName := os.Getenv("TELEGRAPH_AUTHOR_NAME")
	defaultAuthorURL := os.Getenv("TELEGRAPH_AUTHOR_URL")

	// Parse command-line flags (environment variables used as defaults)
	shortName := flag.String("short-name", defaultShortName, "Short name for the Telegraph account (required, can also be set via TELEGRAPH_SHORT_NAME env var)")
	authorName := flag.String("author-name", defaultAuthorName, "Author name (optional, can also be set via TELEGRAPH_AUTHOR_NAME env var)")
	authorURL := flag.String("author-url", defaultAuthorURL, "Author URL (optional, can also be set via TELEGRAPH_AUTHOR_URL env var)")
	flag.Parse()

	// Validate required flag
	if *shortName == "" {
		fmt.Fprintf(os.Stderr, "Error: -short-name is required (or set TELEGRAPH_SHORT_NAME environment variable)\n")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	// Create Telegraph client
	client := telegraph.NewDefaultClient()

	// Build request
	req := telegraph.CreateAccountRequest{
		ShortName: *shortName,
	}

	if *authorName != "" {
		req.AuthorName = authorName
	}

	if *authorURL != "" {
		req.AuthorURL = authorURL
	}

	// Create account
	account, err := client.CreateAccount(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Output the access token
	if account.AccessToken != nil {
		fmt.Println(*account.AccessToken)
	} else {
		fmt.Fprintf(os.Stderr, "Error: no access token received\n")
		os.Exit(1)
	}
}
