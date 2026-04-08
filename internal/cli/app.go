package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bssm-oss/globalAI/internal/browser"
	installer "github.com/bssm-oss/globalAI/internal/install"
	"github.com/bssm-oss/globalAI/internal/source"
	"github.com/bssm-oss/globalAI/internal/viewer"
)

type RuntimeEnvironment struct {
	GetWorkingDirectory func() (string, error)
	GetHomeDirectory    func() (string, error)
	GetExecutablePath   func() (string, error)
	LookupEnv           func(string) (string, bool)
	InstallSelf         func(installer.Config) (installer.Result, error)
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
	if env.GetExecutablePath == nil {
		env.GetExecutablePath = os.Executable
	}
	if env.LookupEnv == nil {
		env.LookupEnv = os.LookupEnv
	}
	if env.InstallSelf == nil {
		env.InstallSelf = installer.Install
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
	case "install":
		return a.runInstall(args[1:], stdout)
	case "web":
		return a.runWeb(ctx, args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText)
	}
}

func (a *App) runInstall(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	destFlag := fs.String("dest", "", "Destination directory for the installed binary.")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, installUsageText)
			return nil
		}
		return err
	}
	if len(fs.Args()) > 0 {
		return errors.New("install does not accept positional arguments")
	}

	home, err := a.env.GetHomeDirectory()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}
	executable, err := a.env.GetExecutablePath()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}
	pathValue, _ := a.env.LookupEnv("PATH")
	goBin, _ := a.env.LookupEnv("GOBIN")
	goPath, _ := a.env.LookupEnv("GOPATH")
	shellPath, _ := a.env.LookupEnv("SHELL")

	result, err := a.env.InstallSelf(installer.Config{
		ExecutablePath: executable,
		HomeDirectory:  home,
		PathValue:      pathValue,
		ShellPath:      shellPath,
		GOBIN:          goBin,
		GOPATH:         goPath,
		Destination:    *destFlag,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "installed globalai to %s\n", result.BinaryPath)
	if result.OnPath {
		fmt.Fprintln(stdout, "globalai is ready on PATH")
		fmt.Fprintln(stdout, "run: globalai --help")
		return nil
	}
	if result.ShellConfigPath != "" {
		if result.ShellConfigUpdated {
			fmt.Fprintf(stdout, "updated %s to include %s\n", result.ShellConfigPath, result.Directory)
		} else {
			fmt.Fprintf(stdout, "%s already includes setup for %s\n", result.ShellConfigPath, result.Directory)
		}
		if result.ReloadCommand != "" {
			fmt.Fprintf(stdout, "run: %s\n", result.ReloadCommand)
		}
	}
	fmt.Fprintf(stdout, "run now: %s --help\n", result.BinaryPath)
	return nil
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
  globalai install [--dest DIR]

Commands:
  web      Start a local viewer for curated prompt and config sources.
  install  Install the current globalai binary into a user bin directory.
  help     Show this help text.
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

const installUsageText = `globalai install copies the current globalai binary into a user bin directory.

Usage:
  globalai install [--dest DIR]

Flags:
  --dest  Destination directory for the installed binary.
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
