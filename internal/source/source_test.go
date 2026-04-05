package source

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverLoadsProjectAndGlobalSources(t *testing.T) {
	root := t.TempDir()
	home := t.TempDir()

	writeFile(t, filepath.Join(root, "AGENTS.md"), "project agent")
	writeFile(t, filepath.Join(root, ".claude", "local.md"), "local prompt")
	writeFile(t, filepath.Join(home, ".claude", "global.md"), "global prompt")
	writeFile(t, filepath.Join(home, "AGENTS.md"), "home agent")
	writeFile(t, filepath.Join(root, ".cursor", "rules", "rule.mdc"), "cursor rule")
	writeFile(t, filepath.Join(root, "ignored.bin"), "ignored")

	result, err := Discover(Config{RootDirectory: root, HomeDirectory: home})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(result.Sources) != 5 {
		t.Fatalf("expected 5 sources, got %d", len(result.Sources))
	}
	if len(result.Families) != 2 {
		t.Fatalf("expected 2 families, got %d", len(result.Families))
	}
}

func TestDiscoverRejectsFileRoot(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "AGENTS.md")
	writeFile(t, path, "content")

	_, err := Discover(Config{RootDirectory: path})
	if err == nil {
		t.Fatal("expected error for file root")
	}
}

func TestDiscoverRejectsSymlinkedFile(t *testing.T) {
	root := t.TempDir()
	targetDir := t.TempDir()
	target := filepath.Join(targetDir, "secret.md")
	writeFile(t, target, "outside")
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	linkPath := filepath.Join(root, ".claude", "linked.md")
	if err := os.Symlink(target, linkPath); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	result, err := Discover(Config{RootDirectory: root, HomeDirectory: t.TempDir()})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(result.Sources) != 0 {
		t.Fatalf("expected symlinked file to be skipped, got %d sources", len(result.Sources))
	}
}

func TestDiscoverRejectsLargeFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AGENTS.md"), strings.Repeat("a", maxFileSizeBytes+1))

	result, err := Discover(Config{RootDirectory: root, HomeDirectory: t.TempDir()})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	if len(result.Sources) != 0 {
		t.Fatalf("expected oversized file to be skipped, got %d sources", len(result.Sources))
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
