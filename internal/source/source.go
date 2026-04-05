package source

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const maxFileSizeBytes = 1024 * 1024

type Config struct {
	RootDirectory string
	HomeDirectory string
}

type Result struct {
	GeneratedAt   time.Time `json:"generatedAt"`
	RootDirectory string    `json:"root"`
	TotalSources  int       `json:"totalSources"`
	Families      []Family  `json:"families"`
	Sources       []Source  `json:"sources"`
}

type Family struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Source struct {
	ID           string    `json:"id"`
	Family       string    `json:"family"`
	Label        string    `json:"label"`
	RelativePath string    `json:"relativePath"`
	Kind         string    `json:"kind"`
	SizeBytes    int64     `json:"sizeBytes"`
	ModifiedAt   time.Time `json:"modifiedAt"`
	Content      string    `json:"content"`
}

type familySpec struct {
	name      string
	baseDir   string
	fileNames []string
	dirNames  []string
}

func DefaultWorkingDirectory() (string, error) {
	return os.Getwd()
}

func DefaultHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

func Discover(cfg Config) (*Result, error) {
	if cfg.RootDirectory == "" {
		return nil, errors.New("root directory is required")
	}
	root, err := filepath.Abs(cfg.RootDirectory)
	if err != nil {
		return nil, fmt.Errorf("resolve root directory: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat root directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root path %q is not a directory", root)
	}

	home := cfg.HomeDirectory
	if home != "" {
		home, err = filepath.Abs(home)
		if err != nil {
			return nil, fmt.Errorf("resolve home directory: %w", err)
		}
	}

	result := &Result{
		GeneratedAt:   time.Now().UTC(),
		RootDirectory: root,
		Families:      []Family{},
		Sources:       []Source{},
	}
	seen := map[string]struct{}{}

	for _, spec := range familySpecs(root, home) {
		count := 0
		for _, name := range spec.fileNames {
			count += appendFile(result, seen, spec.name, spec.baseDir, filepath.Join(spec.baseDir, name), root)
		}
		for _, dirName := range spec.dirNames {
			count += appendDirectory(result, seen, spec.name, spec.baseDir, filepath.Join(spec.baseDir, dirName), root)
		}
		if count > 0 {
			result.Families = append(result.Families, Family{Name: spec.name, Count: count})
		}
	}

	sort.Slice(result.Families, func(i, j int) bool {
		return result.Families[i].Name < result.Families[j].Name
	})
	sort.Slice(result.Sources, func(i, j int) bool {
		if result.Sources[i].Family == result.Sources[j].Family {
			return result.Sources[i].RelativePath < result.Sources[j].RelativePath
		}
		return result.Sources[i].Family < result.Sources[j].Family
	})
	result.TotalSources = len(result.Sources)

	return result, nil
}

func familySpecs(root string, home string) []familySpec {
	specs := []familySpec{{
		name:      "project",
		baseDir:   root,
		fileNames: []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md", filepath.Join(".github", "copilot-instructions.md")},
		dirNames:  []string{".claude", filepath.Join(".cursor", "rules"), ".sisyphus"},
	}}
	if home != "" {
		specs = append(specs, familySpec{
			name:      "global",
			baseDir:   home,
			fileNames: []string{"AGENTS.md"},
			dirNames:  []string{".claude", filepath.Join(".cursor", "rules"), ".sisyphus"},
		})
	}
	return specs
}

func appendFile(result *Result, seen map[string]struct{}, family string, familyBase string, path string, root string) int {
	if isSymlink(path) {
		return 0
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return 0
	}
	if !isAllowedFile(path, info.Size()) {
		return 0
	}
	if appendSource(result, seen, family, familyBase, path, root, info) {
		return 1
	}
	return 0
}

func appendDirectory(result *Result, seen map[string]struct{}, family string, familyBase string, dir string, root string) int {
	if isSymlink(dir) {
		return 0
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return 0
	}

	count := 0
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !isAllowedFile(path, info.Size()) {
			return nil
		}
		if appendSource(result, seen, family, familyBase, path, root, info) {
			count++
		}
		return nil
	})
	return count
}

func appendSource(result *Result, seen map[string]struct{}, family string, familyBase string, path string, root string, info os.FileInfo) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	resolvedPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return false
	}
	if !isWithinResolvedBase(resolvedPath, familyBase) {
		return false
	}
	if _, ok := seen[absPath]; ok {
		return false
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		return false
	}
	seen[absPath] = struct{}{}
	result.Sources = append(result.Sources, Source{
		ID:           sourceID(family, absPath),
		Family:       family,
		Label:        filepath.Base(absPath),
		RelativePath: relativePath(absPath, familyBase, root),
		Kind:         strings.TrimPrefix(strings.ToLower(filepath.Ext(absPath)), "."),
		SizeBytes:    info.Size(),
		ModifiedAt:   info.ModTime().UTC(),
		Content:      string(content),
	})
	return true
}

func relativePath(path string, familyBase string, root string) string {
	if rel, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	if rel, err := filepath.Rel(familyBase, path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(filepath.Base(path))
}

func sourceID(family string, path string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ".", "-")
	return family + "-" + replacer.Replace(path)
}

func isAllowedFile(path string, size int64) bool {
	if size > maxFileSizeBytes {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".mdc", ".txt", ".json", ".yaml", ".yml":
		return true
	default:
		return false
	}
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func isWithinBase(path string, base string) bool {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "")
}

func isWithinResolvedBase(path string, base string) bool {
	resolvedBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return false
	}
	return isWithinBase(path, resolvedBase)
}
