package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/bssm-oss/globalAI/internal/browser"
	"github.com/bssm-oss/globalAI/internal/source"
	"github.com/bssm-oss/globalAI/internal/viewer"
)

type RuntimeEnvironment struct {
	GetWorkingDirectory func() (string, error)
	GetHomeDirectory    func() (string, error)
	OpenURL             func(string) error
	DiscoverNow         func(source.Config) (*source.Result, error)
	ServeViewer         func(context.Context, viewer.Config) (*viewer.Session, error)
}

type App struct {
	env RuntimeEnvironment
}

func New(env RuntimeEnvironment) *App {
	if env.GetWorkingDirectory == nil {
		env.GetWorkingDirectory = source.DefaultWorkingDirectory
	}
	if env.GetHomeDirectory == nil {
		env.GetHomeDirectory = source.DefaultHomeDirectory
	}
	if env.OpenURL == nil {
		env.OpenURL = browser.Open
	}
	if env.DiscoverNow == nil {
		env.DiscoverNow = source.Discover
	}
	if env.ServeViewer == nil {
		env.ServeViewer = viewer.Start
	}

	return &App{env: env}
}

func (a *App) Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		a.printUsage(stdout)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printUsage(stdout)
		return nil
	case "web":
		return a.runWeb(ctx, args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText)
	}
}

func (a *App) runWeb(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("web", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	rootFlag := fs.String("root", "", "Root directory used for repo-local discovery.")
	addrFlag := fs.String("addr", "127.0.0.1:0", "Listening address for the local web server.")
	openFlag := newTrackedBool(true)
	noOpenFlag := newTrackedBool(false)
	fs.Var(openFlag, "open", "Open the viewer in a browser after the server starts.")
	fs.Var(noOpenFlag, "no-open", "Do not open the viewer automatically.")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, webUsageText)
			return nil
		}
		return err
	}
	if openFlag.wasSet && noOpenFlag.wasSet {
		return errors.New("--open and --no-open cannot be used together")
	}

	root, err := a.resolveRoot(*rootFlag, fs.Args())
	if err != nil {
		return err
	}
	home, err := a.env.GetHomeDirectory()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}

	result, err := a.env.DiscoverNow(source.Config{RootDirectory: root, HomeDirectory: home})
	if err != nil {
		return err
	}

	session, err := a.env.ServeViewer(ctx, viewer.Config{
		Address:        *addrFlag,
		OpenInBrowser:  !noOpenFlag.value && openFlag.value,
		OpenBrowserURL: a.env.OpenURL,
		Result:         result,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "globalai viewer ready at %s\n", session.URL)
	fmt.Fprintf(stdout, "sources: %d (%s)\n", len(result.Sources), result.RootDirectory)
	if session.BrowserOpenError != nil {
		fmt.Fprintf(stderr, "browser open skipped: %v\n", session.BrowserOpenError)
	}

	<-session.Done
	return session.Err
}

func (a *App) resolveRoot(rootFlag string, trailing []string) (string, error) {
	if len(trailing) > 1 {
		return "", errors.New("web accepts at most one positional path")
	}

	root := rootFlag
	if root == "" && len(trailing) == 1 {
		root = trailing[0]
	}
	if root == "" {
		cwd, err := a.env.GetWorkingDirectory()
		if err != nil {
			return "", fmt.Errorf("resolve working directory: %w", err)
		}
		root = cwd
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve absolute root path: %w", err)
	}
	return absRoot, nil
}

func (a *App) printUsage(w io.Writer) {
	fmt.Fprint(w, usageText)
}

const usageText = `globalai helps inspect AI prompt and instruction sources.

Usage:
  globalai web [path] [--root PATH] [--addr HOST:PORT] [--open|--no-open]

Commands:
  web   Start a local viewer for curated prompt and config sources.
  help  Show this help text.
`

const webUsageText = `globalai web starts a local viewer for curated prompt and config sources.

Usage:
  globalai web [path] [--root PATH] [--addr LOOPBACK_HOST:PORT] [--open|--no-open]

Flags:
  --root     Root directory used for repo-local discovery.
  --addr     Loopback listening address, default 127.0.0.1:0.
  --open     Open the viewer in a browser after startup.
  --no-open  Skip browser opening.
`

type trackedBool struct {
	value  bool
	wasSet bool
}

func newTrackedBool(defaultValue bool) *trackedBool {
	return &trackedBool{value: defaultValue}
}

func (b *trackedBool) String() string {
	return strconv.FormatBool(b.value)
}

func (b *trackedBool) Set(value string) error {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	b.value = parsed
	b.wasSet = true
	return nil
}

func (b *trackedBool) IsBoolFlag() bool {
	return true
}
