package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(url string) error {
	command, args, err := commandFor(url)
	if err != nil {
		return err
	}

	if err := exec.Command(command, args...).Start(); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	return nil
}

func commandFor(url string) (string, []string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "open", []string{url}, nil
	case "linux":
		return "xdg-open", []string{url}, nil
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", url}, nil
	default:
		return "", nil, fmt.Errorf("unsupported operating system %q", runtime.GOOS)
	}
}
