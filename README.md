[![Go Doc](https://pkg.go.dev/badge/github.com/rizalgowandy/mailhog-go?status.svg)](https://pkg.go.dev/github.com/rizalgowandy/mailhog-go?tab=doc)
[![Release](https://img.shields.io/github/release/rizalgowandy/mailhog-go.svg?style=flat-square)](https://github.com/rizalgowandy/mailhog-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/rizalgowandy/mailhog-go)](https://goreportcard.com/report/github.com/rizalgowandy/mailhog-go)
[![Build Status](https://github.com/rizalgowandy/mailhog-go/workflows/Go/badge.svg?branch=main)](https://github.com/rizalgowandy/mailhog-go/actions?query=branch%3Amain)
[![Sourcegraph](https://sourcegraph.com/github.com/rizalgowandy/mailhog-go/-/badge.svg)](https://sourcegraph.com/github.com/rizalgowandy/mailhog-go?badge)

![logo](.github/mailhog-go.png)

## Getting Started

Interact with MailHog API.

## Installation

```shell
go get -v github.com/rizalgowandy/mailhog-go
```

## Quick Start

```go
package main

import (
	"log"
	"context"

	"github.com/rizalgowandy/mailhog-go"
	"github.com/rizalgowandy/mailhog-go/pkg/api"
)

func main() {
	cfg := api.Config{
		HostURL: MailHogContainer.UIEndpoint,
	}
	client, err := mailhog.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	messages, err := client.GetAllMessages(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

For more example, check [here](test/smtp_test.go).

## Supported API

- GET /api/v2/messages
- GET /api/v1/messages/{id}
- DELETE /api/v1/messages
- DELETE /api/v1/messages/{id}
- GET /api/v2/search
