# Assistant

An AI-powered content generation service that uses Google Gemini to create structured documents and publish them to Telegraph.

## Overview

This service provides an HTTP API that accepts content generation requests. Requests are processed asynchronously by Cloud Run jobs that:

1. Generate content using Google Gemini AI
2. Structure the content into documents
3. Publish documents to Telegraph
4. Send email notifications

## Architecture

- **Service** (`cmd/service`): HTTP server that receives requests and queues Cloud Run jobs
- **Job** (`cmd/job`): Worker that processes content generation requests
- **Generators**: Pluggable content generators registered via the registry system

## Building

```bash
make build
```

This builds three binaries:
- `bin/services` - HTTP service
- `bin/job` - Content generation worker
- `bin/tokens` - Telegraph token generator

## Usage

The service expects:
- `X-API-Token` header for authentication
- `X-Content-Type` header specifying which generator to use
- Optional `X-Config-*` headers for generator-specific configuration
- JSON request body with generator-specific payload

## Writing Custom Content Generators

Content generators implement the `ContentGenerator` interface and are registered via the generator registry.

### 1. Implement the Interface

```go
package mygenerator

import (
    "context"
    "github.com/schraf/assistant/pkg/models"
)

type MyGenerator struct {
    // Your generator-specific fields
}

func (g *MyGenerator) Generate(ctx context.Context, request models.ContentRequest, assistant models.Assistant) (*models.Document, error) {
    // Use assistant.Ask() or assistant.StructuredAsk() to generate content
    // Parse request.Body to get your input parameters
    // Return a models.Document with Title, Author, and Sections
}
```

### 2. Create a Factory Function

```go
func NewGenerator(config generators.Config) (models.ContentGenerator, error) {
    // Extract configuration from config map
    // Initialize and return your generator
    return &MyGenerator{
        // ... initialized fields
    }, nil
}
```

### 3. Register the Generator

Register your generator in an `init()` function (typically in a separate package that gets imported):

```go
package mygenerator

import (
    "github.com/schraf/assistant/pkg/generators"
)

func init() {
    generators.MustRegister("my-generator", NewGenerator)
}
```

### 4. Import the Package

Import your generator package in `cmd/job/main.go` to ensure registration:

```go
import (
    _ "path/to/mygenerator"
)
```

### Example

The `ContentGenerator` interface:

```go
type ContentGenerator interface {
    Generate(ctx context.Context, request ContentRequest, assistant Assistant) (*Document, error)
}
```

The `Assistant` interface provides:
- `Ask(ctx, persona, request) (*string, error)` - Generate text responses
- `StructuredAsk(ctx, persona, request, schema) (json.RawMessage, error)` - Generate structured JSON responses

The `Document` model:

```go
type Document struct {
    Title    string
    Author   string
    Sections []DocumentSection
}

type DocumentSection struct {
    Title      string
    Paragraphs []string
}
```

Configuration is passed via the `Config` map (from `X-Config-*` headers) and request data via `ContentRequest.Body`.

## Deployment

The project includes Terraform configurations for deploying to Google Cloud Platform. See `terraform/` directory and the Makefile for deployment commands.

## License

See LICENSE file for details.
