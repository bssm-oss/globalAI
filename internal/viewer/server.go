package viewer

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bssm-oss/globalAI/internal/source"
)

//go:embed static/*
var embeddedAssets embed.FS

type Config struct {
	Address        string
	OpenInBrowser  bool
	OpenBrowserURL func(string) error
	Result         *source.Result
}

type Session struct {
	URL              string
	Done             <-chan struct{}
	Err              error
	BrowserOpenError error
}

func Start(ctx context.Context, cfg Config) (*Session, error) {
	if cfg.Result == nil {
		return nil, fmt.Errorf("viewer result is required")
	}
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:0"
	}
	if err := validateLoopbackAddress(cfg.Address); err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	assetFS, err := fs.Sub(embeddedAssets, "static")
	if err != nil {
		return nil, fmt.Errorf("load embedded assets: %w", err)
	}

	mux := http.NewServeMux()
	payload, err := json.Marshal(cfg.Result)
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("marshal viewer payload: %w", err)
	}
	mux.HandleFunc("/api/sources", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(payload)
	})
	mux.Handle("/", http.FileServer(http.FS(assetFS)))

	server := &http.Server{Handler: mux}
	done := make(chan struct{})
	session := &Session{URL: fmt.Sprintf("http://%s", listener.Addr().String()), Done: done}

	var once sync.Once
	finish := func(err error) {
		once.Do(func() {
			session.Err = err
			close(done)
		})
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	go func() {
		defer listener.Close()
		if serveErr := server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			finish(serveErr)
			return
		}
		finish(nil)
	}()

	if cfg.OpenInBrowser && cfg.OpenBrowserURL != nil {
		session.BrowserOpenError = cfg.OpenBrowserURL(session.URL)
	}

	return session, nil
}

const shutdownTimeout = 5 * time.Second

func validateLoopbackAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid address %q: %w", address, err)
	}
	host = strings.Trim(host, "[]")
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("address %q must bind to a loopback host", address)
	}
	return nil
}
