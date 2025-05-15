package mailhog

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rizalgowandy/mailhog-go/pkg/api"
	"github.com/rizalgowandy/mailhog-go/pkg/entity"
)

// NewClient creates a client to interact with XYZ API.
func NewClient(cfg api.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Client{
		cli: api.NewRestyClient(cfg),
	}, nil
}

// Client is the main client to interact with XYZ API.
type Client struct {
	cli *resty.Client
}

// GetAllMessages retrieves all messages
func (c *Client) GetAllMessages(ctx context.Context) ([]entity.Message, error) {
	url := "/api/v2/messages"

	var (
		content struct {
			Items []entity.Message `json:"items"`
		}
		contentErr entity.ErrResp
	)

	_, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&content).
		SetError(&contentErr).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("get all messages: %w", err)
	}

	return content.Items, nil
}

// GetMessage retrieves a specific message by ID
func (c *Client) GetMessage(ctx context.Context, id string) (*entity.Message, error) {
	url := "/api/v1/messages/{id}"

	var (
		content    entity.Message
		contentErr entity.ErrResp
	)

	_, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetPathParam("id", id).
		SetResult(&content).
		SetError(&contentErr).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("get message with ID %s: %w", id, err)
	}
	if content.Content.Size == 0 {
		return nil, fmt.Errorf("message with ID %s not found", id)
	}

	return &content, nil
}

// DeleteAllMessages deletes all messages
func (c *Client) DeleteAllMessages(ctx context.Context) error {
	url := "/api/v1/messages"

	var contentErr entity.ErrResp

	_, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetError(&contentErr).
		Delete(url)
	if err != nil {
		return fmt.Errorf("delete all messages: %w", err)
	}

	return nil
}

// DeleteMessage deletes a specific message by ID
func (c *Client) DeleteMessage(ctx context.Context, id string) error {
	url := "/api/v1/messages/{id}"

	var contentErr entity.ErrResp

	_, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetPathParam("id", id).
		SetError(&contentErr).
		Delete(url)
	if err != nil {
		return fmt.Errorf("delete message with ID %s: %w", id, err)
	}

	return nil
}

// SearchMessages searches for messages based on the provided parameters
func (c *Client) SearchMessages(ctx context.Context, kind, query string, start, limit int) ([]entity.Message, error) {
	url := "/api/v2/search"

	var (
		content struct {
			Items []entity.Message `json:"items"`
			Total int              `json:"total"`
		}
		contentErr entity.ErrResp
	)

	_, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParam("kind", kind).
		SetQueryParam("query", query).
		SetQueryParam("start", fmt.Sprintf("%d", start)).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&content).
		SetError(&contentErr).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("search messages with kind %s and query %s: %w", kind, query, err)
	}

	return content.Items, nil
}

// LatestFrom retrieves the latest message from a specific sender
func (c *Client) LatestFrom(ctx context.Context, from string) (*entity.Message, error) {
	messages, err := c.SearchMessages(ctx, "from", from, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("get latest message from %s: %w", from, err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[0], nil
}

// LatestTo retrieves the latest message to a specific recipient
func (c *Client) LatestTo(ctx context.Context, to string) (*entity.Message, error) {
	messages, err := c.SearchMessages(ctx, "to", to, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("get latest message to %s: %w", to, err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[0], nil
}

// LatestContaining retrieves the latest message containing specific content
func (c *Client) LatestContaining(ctx context.Context, query string) (*entity.Message, error) {
	messages, err := c.SearchMessages(ctx, "containing", query, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("get latest message containing %s: %w", query, err)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return &messages[0], nil
}
