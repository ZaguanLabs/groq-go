# Groq Go SDK

A Go client library for accessing the Groq API.

## Installation

```bash
go get github.com/ZaguanLabs/groq-go/groq
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZaguanLabs/groq-go/groq"
)

func main() {
	client, err := groq.NewClient(
		groq.WithAPIKey(os.Getenv("GROQ_API_KEY")),
	)
	if err != nil {
		panic(err)
	}

	// Use client to make requests...
    _ = client
}
```

## Project Structure

- `groq/`: Main SDK source code
- `docs/`: Documentation and plans
