package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	ExecutablePath string
	HomeDirectory  string
	PathValue      string
	ShellPath      string
	GOBIN          string
	GOPATH         string
	Destination    string
}

type Result struct {
	BinaryPath         string
	Directory          string
	OnPath             bool
	ShellConfigPath    string
	ShellConfigUpdated bool
	NeedsShellReload   bool
	ReloadCommand      string
}

func Install(cfg Config) (Result, error) {
	if cfg.ExecutablePath == "" {
		return Result{}, fmt.Errorf("resolve executable path: empty path")
	}
	if cfg.HomeDirectory == "" {
		return Result{}, fmt.Errorf("resolve home directory: empty path")
	}

	destination, err := chooseDestinationDir(cfg)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return Result{}, fmt.Errorf("create install directory: %w", err)
	}

	binaryPath := filepath.Join(destination, binaryName())
	if err := copyExecutable(cfg.ExecutablePath, binaryPath); err != nil {
		return Result{}, err
	}

	result := Result{
		BinaryPath: binaryPath,
		Directory:  destination,
		OnPath:     pathContainsDir(cfg.PathValue, destination),
	}
	if result.OnPath {
		return result, nil
	}

	rcPath, line, reloadCommand := shellSetup(cfg.ShellPath, cfg.HomeDirectory, destination)
	if rcPath == "" {
		return result, nil
	}
	updated, err := appendIfMissing(rcPath, line)
	if err != nil {
		return result, err
	}
	result.ShellConfigPath = rcPath
	result.ShellConfigUpdated = updated
	result.NeedsShellReload = true
	result.ReloadCommand = reloadCommand
	return result, nil
}

func chooseDestinationDir(cfg Config) (string, error) {
	if cfg.Destination != "" {
		return filepath.Abs(cfg.Destination)
	}
	if pathDir, ok := firstWritableHomePathDir(cfg.PathValue, cfg.HomeDirectory); ok {
		return pathDir, nil
	}
	if cfg.GOBIN != "" {
		return cfg.GOBIN, nil
	}
	if cfg.GOPATH != "" {
		return filepath.Join(cfg.GOPATH, "bin"), nil
	}
	return filepath.Join(cfg.HomeDirectory, "go", "bin"), nil
}

func firstWritableHomePathDir(pathValue, homeDirectory string) (string, bool) {
	for _, entry := range filepath.SplitList(pathValue) {
		if entry == "" {
			continue
		}
		if !isWithinHome(entry, homeDirectory) {
			continue
		}
		if isWritableDir(entry) {
			return entry, true
		}
	}
	return "", false
}

func isWithinHome(path, homeDirectory string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absHome, err := filepath.Abs(homeDirectory)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absHome, absPath)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "")
}

func isWritableDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	file, err := os.CreateTemp(path, ".globalai-write-check-")
	if err != nil {
		return false
	}
	name := file.Name()
	_ = file.Close()
	_ = os.Remove(name)
	return true
}

func copyExecutable(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open current executable: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create binary directory: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(dst), ".globalai-install-")
	if err != nil {
		return fmt.Errorf("create temp binary: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := io.Copy(temp, in); err != nil {
		_ = temp.Close()
		return fmt.Errorf("copy binary: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temp binary: %w", err)
	}
	if err := os.Chmod(tempPath, 0o755); err != nil {
		return fmt.Errorf("mark binary executable: %w", err)
	}
	if err := os.Rename(tempPath, dst); err != nil {
		return fmt.Errorf("install binary: %w", err)
	}
	return nil
}

func pathContainsDir(pathValue, dir string) bool {
	target := filepath.Clean(dir)
	for _, entry := range filepath.SplitList(pathValue) {
		if entry != "" && filepath.Clean(entry) == target {
			return true
		}
	}
	return false
}

func shellSetup(shellPath, homeDirectory, installDir string) (string, string, string) {
	switch filepath.Base(shellPath) {
	case "zsh":
		rc := filepath.Join(homeDirectory, ".zshrc")
		return rc, fmt.Sprintf("export PATH=%q:$PATH", installDir), fmt.Sprintf("source %s", rc)
	case "bash":
		rc := filepath.Join(homeDirectory, ".bashrc")
		return rc, fmt.Sprintf("export PATH=%q:$PATH", installDir), fmt.Sprintf("source %s", rc)
	case "fish":
		rc := filepath.Join(homeDirectory, ".config", "fish", "config.fish")
		return rc, fmt.Sprintf("fish_add_path %q", installDir), fmt.Sprintf("source %s", rc)
	default:
		return "", "", ""
	}
}

func appendIfMissing(path, line string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read shell config: %w", err)
	}
	if strings.Contains(string(content), line) {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("create shell config directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return false, fmt.Errorf("open shell config: %w", err)
	}
	defer file.Close()

	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		if _, err := file.WriteString("\n"); err != nil {
			return false, fmt.Errorf("write shell config newline: %w", err)
		}
	}
	if _, err := file.WriteString(line + "\n"); err != nil {
		return false, fmt.Errorf("write shell config: %w", err)
	}
	return true, nil
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "globalai.exe"
	}
	return "globalai"
}
