package test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var MailHogContainer *MailHogTestContainer

// MailHogTestContainer represents a MailHog container with helper methods
type MailHogTestContainer struct {
	testcontainers.Container
	DSN        string
	SMTPPort   string
	UIEndpoint string
}

// SetupMailHog starts a MailHog container and sets up its DSN.
func SetupMailHog(ctx context.Context) error {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mailhog/mailhog:v1.0.1",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("1025/tcp"), // SMTP port
				wait.ForListeningPort("8025/tcp"), // API/Web UI port
				wait.ForHTTP("/api/v1/messages").
					WithPort("8025/tcp").
					WithStatusCodeMatcher(func(status int) bool {
						return status == http.StatusOK
					}),
			).WithDeadline(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("failed to start MailHog container: %w", err)
	}

	MailHogContainer = &MailHogTestContainer{
		Container: container,
	}

	smtpPort, err := MailHogContainer.MappedPort(ctx, "1025")
	if err != nil {
		return fmt.Errorf("failed to get SMTP port: %w", err)
	}
	MailHogContainer.SMTPPort = smtpPort.Port()

	httpPort, err := MailHogContainer.MappedPort(ctx, "8025")
	if err != nil {
		return fmt.Errorf("failed to get HTTP port: %w", err)
	}

	MailHogContainer.DSN = fmt.Sprintf("smtp://localhost:%s", smtpPort.Port())
	MailHogContainer.UIEndpoint = fmt.Sprintf("http://localhost:%s", httpPort.Port())
	return nil
}
