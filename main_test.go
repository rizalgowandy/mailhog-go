package mailhog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rizalgowandy/mailhog-go/pkg/api"
	"github.com/rizalgowandy/mailhog-go/pkg/entity"
	"github.com/stretchr/testify/require"
)

func setupClientServer(t *testing.T) (*Client, func()) {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v2/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"ID": "msg-1",
				"Content": map[string]any{
					"Size": 10,
				},
			}},
		})
	})

	mux.HandleFunc("/api/v1/messages/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/messages/")

		switch r.Method {
		case http.MethodGet:
			if id == "missing" {
				w.Header().Set("Content-Type", "application/json")

				_ = json.NewEncoder(w).Encode(map[string]any{
					"ID": "missing",
					"Content": map[string]any{
						"Size": 0,
					},
				})

				return
			}

			w.Header().Set("Content-Type", "application/json")

			_ = json.NewEncoder(w).Encode(map[string]any{
				"ID": id,
				"Content": map[string]any{
					"Size": 7,
				},
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/v2/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		kind := q.Get("kind")
		query := q.Get("query")
		start, _ := strconv.Atoi(q.Get("start"))
		limit, _ := strconv.Atoi(q.Get("limit"))

		if start < 0 || limit <= 0 {
			http.Error(w, "invalid pagination", http.StatusBadRequest)
			return
		}

		if query == "none" {
			w.Header().Set("Content-Type", "application/json")

			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{},
				"total": 0,
			})

			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"ID": fmt.Sprintf("%s-%s", kind, query),
				"Content": map[string]any{
					"Size": 1,
				},
			}},
			"total": 1,
		})
	})

	srv := httptest.NewServer(mux)

	client, err := NewClient(api.Config{HostURL: srv.URL, Timeout: time.Second})
	require.NoError(t, err)

	return client, srv.Close
}

func TestNewClientAppliesDefaults(t *testing.T) {
	client, err := NewClient(api.Config{})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.cli)
	require.Equal(t, "http://localhost:8025", client.cli.BaseURL)
}

func TestClientMessageLifecycleAndSearch(t *testing.T) {
	client, cleanup := setupClientServer(t)
	defer cleanup()

	ctx := context.Background()

	all, err := client.GetAllMessages(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, "msg-1", all[0].ID)

	msg, err := client.GetMessage(ctx, "abc")
	require.NoError(t, err)
	require.Equal(t, "abc", msg.ID)

	msg, err = client.GetMessage(ctx, "missing")
	require.Error(t, err)
	require.Nil(t, msg)
	require.Contains(t, err.Error(), "not found")

	err = client.DeleteAllMessages(ctx)
	require.NoError(t, err)

	err = client.DeleteMessage(ctx, "abc")
	require.NoError(t, err)

	found, err := client.SearchMessages(ctx, "from", "sender@example.com", 0, 1)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Equal(t, "from-sender@example.com", found[0].ID)
}

func TestClientLatestHelpers(t *testing.T) {
	client, cleanup := setupClientServer(t)
	defer cleanup()

	ctx := context.Background()

	from, err := client.LatestFrom(ctx, "alice@example.com")
	require.NoError(t, err)
	require.Equal(t, "from-alice@example.com", from.ID)

	to, err := client.LatestTo(ctx, "bob@example.com")
	require.NoError(t, err)
	require.Equal(t, "to-bob@example.com", to.ID)

	containing, err := client.LatestContaining(ctx, "invoice")
	require.NoError(t, err)
	require.Equal(t, "containing-invoice", containing.ID)
}

func TestClientLatestHelpersNoMessages(t *testing.T) {
	client, cleanup := setupClientServer(t)
	defer cleanup()

	ctx := context.Background()

	msg, err := client.LatestFrom(ctx, "none")
	require.Error(t, err)
	require.Nil(t, msg)

	msg, err = client.LatestTo(ctx, "none")
	require.Error(t, err)
	require.Nil(t, msg)

	msg, err = client.LatestContaining(ctx, "none")
	require.Error(t, err)
	require.Nil(t, msg)
}

func TestClientRequestErrors(t *testing.T) {
	cfg := api.Config{
		HostURL:          "http://127.0.0.1:1",
		Timeout:          100 * time.Millisecond,
		RetryCount:       0,
		RetryMaxWaitTime: 100 * time.Millisecond,
	}

	client, err := NewClient(cfg)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = client.GetAllMessages(ctx)
	require.Error(t, err)

	_, err = client.GetMessage(ctx, "id")
	require.Error(t, err)

	err = client.DeleteAllMessages(ctx)
	require.Error(t, err)

	err = client.DeleteMessage(ctx, "id")
	require.Error(t, err)

	_, err = client.SearchMessages(ctx, "from", "test", 0, 1)
	require.Error(t, err)

	_, err = client.LatestFrom(ctx, "x")
	require.Error(t, err)

	_, err = client.LatestTo(ctx, "x")
	require.Error(t, err)

	_, err = client.LatestContaining(ctx, "x")
	require.Error(t, err)
}

func TestErrRespTypeCompile(t *testing.T) {
	_ = entity.ErrResp{}
}
