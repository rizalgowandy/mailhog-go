package test

import (
	"context"
	"testing"

	"github.com/rizalgowandy/gdk/pkg/converter"
	"github.com/rizalgowandy/mailhog-go"
	"github.com/rizalgowandy/mailhog-go/pkg/api"
	"github.com/stretchr/testify/require"
	"gopkg.in/mail.v2"
)

func TestSendEmail(t *testing.T) {
	// Step 1: Send an email using gopkg.in/mail.v2
	m := mail.NewMessage()
	m.SetHeader("From", "sender@example.com")
	m.SetHeader("To", "recipient@example.com")
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email.")

	d := mail.NewDialer("localhost", converter.Int(MailHogContainer.SMTPPort), "", "")
	err := d.DialAndSend(m)
	require.NoError(t, err, "Failed to send email")

	// Step 2: Initialize the Client
	cfg := api.Config{
		HostURL: MailHogContainer.UIEndpoint,
	}
	client, err := mailhog.NewClient(cfg)
	require.NoError(t, err, "Failed to create client")

	// Step 3: Get all messages
	ctx := context.Background()
	messages, err := client.GetAllMessages(ctx)
	require.NoError(t, err, "Failed to get all messages")
	require.NotEmpty(t, messages, "No messages found")

	// Step 4: Get a specific message by ID
	messageID := messages[0].ID
	message, err := client.GetMessage(ctx, messageID)
	require.NoError(t, err, "Failed to get message by ID")
	require.NotNil(t, message, "Message not found")
	t.Logf("Message ID: %s", message.ID)

	// Step 5: Delete all messages
	err = client.DeleteAllMessages(ctx)
	require.NoError(t, err, "Failed to delete all messages")

	// Step 6: Verify that all messages are deleted
	messages, err = client.GetAllMessages(ctx)
	require.NoError(t, err, "Failed to get all messages after deletion")
	require.Empty(t, messages, "Messages still exist after deletion")

	// Step 7: Verify that the message is deleted
	message, err = client.GetMessage(ctx, messageID)
	require.Error(t, err)
	require.Nil(t, message, "Expected nil message after deletion")
	require.Contains(t, err.Error(), "not found", "Expected 'message not found' error")
}
