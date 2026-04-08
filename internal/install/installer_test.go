package install

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChooseDestinationDirPrefersWritableHomePath(t *testing.T) {
	home := t.TempDir()
	pathDir := filepath.Join(home, "bin")
	if err := os.MkdirAll(pathDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	got, err := chooseDestinationDir(Config{
		HomeDirectory: home,
		PathValue:     pathDir + string(os.PathListSeparator) + "/usr/bin",
	})
	if err != nil {
		t.Fatalf("chooseDestinationDir() error = %v", err)
	}
	if got != pathDir {
		t.Fatalf("chooseDestinationDir() = %q, want %q", got, pathDir)
	}
}

func TestChooseDestinationDirFallsBackToGoBin(t *testing.T) {
	home := t.TempDir()
	gobin := filepath.Join(home, "custom-bin")

	got, err := chooseDestinationDir(Config{
		HomeDirectory: home,
		PathValue:     "/usr/bin:/bin",
		GOBIN:         gobin,
	})
	if err != nil {
		t.Fatalf("chooseDestinationDir() error = %v", err)
	}
	if got != gobin {
		t.Fatalf("chooseDestinationDir() = %q, want %q", got, gobin)
	}
}

func TestInstallUpdatesShellConfigWhenDestinationNotOnPath(t *testing.T) {
	home := t.TempDir()
	exe := filepath.Join(home, "globalai-current")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := Install(Config{
		ExecutablePath: exe,
		HomeDirectory:  home,
		PathValue:      "/usr/bin:/bin",
		ShellPath:      "/bin/zsh",
		GOBIN:          filepath.Join(home, "go-bin"),
	})
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	if result.OnPath {
		t.Fatal("expected destination not to be on PATH")
	}
	if !result.ShellConfigUpdated {
		t.Fatal("expected shell config to be updated")
	}
	if !result.NeedsShellReload {
		t.Fatal("expected shell reload to be required")
	}
	if _, err := os.Stat(result.BinaryPath); err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(home, ".zshrc"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) == "" {
		t.Fatal("expected zshrc to contain PATH update")
	}
}

func TestInstallUsesExistingWritablePathDirectory(t *testing.T) {
	home := t.TempDir()
	pathDir := filepath.Join(home, "bin")
	if err := os.MkdirAll(pathDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	exe := filepath.Join(home, "globalai-current")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := Install(Config{
		ExecutablePath: exe,
		HomeDirectory:  home,
		PathValue:      pathDir + string(os.PathListSeparator) + "/usr/bin",
		ShellPath:      "/bin/zsh",
	})
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	if !result.OnPath {
		t.Fatal("expected install destination to already be on PATH")
	}
	if result.ShellConfigPath != "" {
		t.Fatalf("expected no shell config update, got %q", result.ShellConfigPath)
	}
	if result.Directory != pathDir {
		t.Fatalf("result.Directory = %q, want %q", result.Directory, pathDir)
	}
}
