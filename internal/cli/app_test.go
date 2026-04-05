package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/bssm-oss/globalAI/internal/source"
	"github.com/bssm-oss/globalAI/internal/viewer"
)

func TestRunHelp(t *testing.T) {
	app := New(RuntimeEnvironment{})
	var stdout bytes.Buffer

	if err := app.Run(context.Background(), nil, &stdout, &bytes.Buffer{}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := stdout.String(); !strings.Contains(got, "globalai helps inspect") {
		t.Fatalf("usage text missing, got %q", got)
	}
}

func TestRunWebHonorsNoOpen(t *testing.T) {
	t.Helper()

	var opened bool
	app := New(RuntimeEnvironment{
		GetWorkingDirectory: func() (string, error) { return "/workspace", nil },
		GetHomeDirectory:    func() (string, error) { return "/home/tester", nil },
		OpenURL: func(string) error {
			opened = true
			return nil
		},
		DiscoverNow: func(cfg source.Config) (*source.Result, error) {
			if cfg.RootDirectory != "/workspace" {
				t.Fatalf("unexpected root directory %q", cfg.RootDirectory)
			}
			return &source.Result{RootDirectory: cfg.RootDirectory}, nil
		},
		ServeViewer: func(_ context.Context, cfg viewer.Config) (*viewer.Session, error) {
			if cfg.OpenInBrowser {
				t.Fatalf("expected browser opening to be disabled")
			}
			done := make(chan struct{})
			close(done)
			return &viewer.Session{URL: "http://127.0.0.1:4000", Done: done}, nil
		},
	})

	var stdout, stderr bytes.Buffer
	if err := app.Run(context.Background(), []string{"web", "--no-open"}, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if opened {
		t.Fatal("expected OpenURL not to be called")
	}
	if got := stdout.String(); !strings.Contains(got, "viewer ready") {
		t.Fatalf("expected startup output, got %q", got)
	}
}

func TestRunWebValidatesConflictingFlags(t *testing.T) {
	app := New(RuntimeEnvironment{})
	err := app.Run(context.Background(), []string{"web", "--open", "--no-open"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("expected conflicting flag error, got %v", err)
	}
}

func TestRunWebReturnsServeError(t *testing.T) {
	serveErr := errors.New("serve failed")
	app := New(RuntimeEnvironment{
		GetWorkingDirectory: func() (string, error) { return "/workspace", nil },
		GetHomeDirectory:    func() (string, error) { return "/home/tester", nil },
		DiscoverNow: func(cfg source.Config) (*source.Result, error) {
			return &source.Result{RootDirectory: cfg.RootDirectory}, nil
		},
		ServeViewer: func(_ context.Context, _ viewer.Config) (*viewer.Session, error) {
			return nil, serveErr
		},
	})

	err := app.Run(context.Background(), []string{"web"}, &bytes.Buffer{}, &bytes.Buffer{})
	if !errors.Is(err, serveErr) {
		t.Fatalf("expected serve error %v, got %v", serveErr, err)
	}
}

func TestRunWebHelp(t *testing.T) {
	app := New(RuntimeEnvironment{})
	var stdout, stderr bytes.Buffer

	if err := app.Run(context.Background(), []string{"web", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got := stdout.String(); !strings.Contains(got, "globalai web starts a local viewer") {
		t.Fatalf("expected web help text, got %q", got)
	}
}
