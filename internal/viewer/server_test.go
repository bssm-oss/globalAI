package viewer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bssm-oss/globalAI/internal/source"
)

func TestStartServesSourcesAPI(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := &source.Result{
		RootDirectory: "/workspace",
		Sources: []source.Source{{
			ID:      "project-agents",
			Content: "hello",
		}},
	}

	session, err := Start(ctx, Config{Address: "127.0.0.1:0", Result: result})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	resp, err := http.Get(session.URL + "/api/sources")
	if err != nil {
		t.Fatalf("GET /api/sources error = %v", err)
	}
	defer resp.Body.Close()

	var payload source.Result
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if len(payload.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(payload.Sources))
	}

	pageResp, err := http.Get(session.URL + "/")
	if err != nil {
		t.Fatalf("GET / error = %v", err)
	}
	defer pageResp.Body.Close()
	body, err := io.ReadAll(pageResp.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !strings.Contains(string(body), "brand-title\">viewer") {
		t.Fatalf("expected embedded viewer shell, got %q", string(body))
	}

	cancel()
	select {
	case <-session.Done:
	case <-time.After(2 * time.Second):
		t.Fatal("viewer did not shut down after context cancellation")
	}
}

func TestStartRejectsNonLoopbackAddress(t *testing.T) {
	_, err := Start(context.Background(), Config{Address: "0.0.0.0:0", Result: &source.Result{}})
	if err == nil || !strings.Contains(err.Error(), "loopback") {
		t.Fatalf("expected loopback validation error, got %v", err)
	}
}
