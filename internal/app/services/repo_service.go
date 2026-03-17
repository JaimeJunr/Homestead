package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DefaultRepoDir is the default directory for the config repo (under user config or home)
const DefaultRepoDirName = "homestead-dotfiles"

// DefaultDotfilesPaths are the paths relative to $HOME to include in the repo (zsh + common)
var DefaultDotfilesPaths = []string{".zshrc", ".zsh"}

// RepoService handles git-based config repo: init, push, clone, restore
type RepoService struct {
	repoDir string // absolute path to repo (e.g. ~/.config/homestead-dotfiles)
}

// NewRepoService creates a repo service. repoDir is the path to the git repo (created on Init or Clone).
func NewRepoService(repoDir string) (*RepoService, error) {
	dir, err := expandHome(repoDir)
	if err != nil {
		return nil, fmt.Errorf("repo dir: %w", err)
	}
	return &RepoService{repoDir: dir}, nil
}

// RepoDir returns the absolute repo directory path
func (rs *RepoService) RepoDir() string {
	return rs.repoDir
}

// IsRepo returns true if repoDir exists and is a git repository
func (rs *RepoService) IsRepo() bool {
	gitDir := filepath.Join(rs.repoDir, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

// InitRepo creates the directory and runs git init
func (rs *RepoService) InitRepo() error {
	if err := os.MkdirAll(rs.repoDir, 0755); err != nil {
		return fmt.Errorf("create repo dir: %w", err)
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = rs.repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// AddRemote adds a remote (e.g. origin, url)
func (rs *RepoService) AddRemote(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = rs.repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git remote add: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// HasRemote returns true if the given remote exists
func (rs *RepoService) HasRemote(name string) bool {
	cmd := exec.Command("git", "remote", "get-url", name)
	cmd.Dir = rs.repoDir
	return cmd.Run() == nil
}

// GetRemoteURL returns the URL for the given remote (e.g. "origin"), or empty string if not set
func (rs *RepoService) GetRemoteURL(name string) string {
	cmd := exec.Command("git", "remote", "get-url", name)
	cmd.Dir = rs.repoDir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CopyToRepo copies paths from homeDir to repoDir. Paths are relative to home (e.g. ".zshrc", ".zsh").
func (rs *RepoService) CopyToRepo(homeDir string, paths []string) error {
	home, err := expandHome(homeDir)
	if err != nil {
		return err
	}
	for _, p := range paths {
		src := filepath.Join(home, p)
		dst := filepath.Join(rs.repoDir, p)
		if err := copyPath(src, dst); err != nil {
			return fmt.Errorf("copy %s: %w", p, err)
		}
	}
	return nil
}

// CommitAll adds all files and commits with the given message
func (rs *RepoService) CommitAll(message string) error {
	add := exec.Command("git", "add", "-A")
	add.Dir = rs.repoDir
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	commit := exec.Command("git", "commit", "-m", message)
	commit.Dir = rs.repoDir
	if out, err := commit.CombinedOutput(); err != nil {
		// Nothing to commit is not an error
		if strings.Contains(string(out), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("git commit: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Push runs git push to the given remote and branch
func (rs *RepoService) Push(remote, branch string) error {
	if branch == "" {
		branch = "main"
	}
	cmd := exec.Command("git", "push", "-u", remote, branch)
	cmd.Dir = rs.repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// CreateGitHubRepoWithGh creates a GitHub repo via gh CLI and pushes the current repoDir.
// Requires gh to be installed and authenticated (gh auth login). private=true creates a private repo.
func CreateGitHubRepoWithGh(repoDir, repoName string, private bool) error {
	visibility := "--public"
	if private {
		visibility = "--private"
	}
	// gh repo create <name> --private|--public --source=. --remote=origin --push
	cmd := exec.Command("gh", "repo", "create", repoName, visibility, "--source", repoDir, "--remote", "origin", "--push")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		if strings.Contains(string(out), "could not find gh") || strings.Contains(string(out), "command not found") {
			return fmt.Errorf("gh não está instalado. Instale: https://cli.github.com — %w", err)
		}
		if strings.Contains(string(out), "authentication required") || strings.Contains(string(out), "failed to authenticate") {
			return fmt.Errorf("gh não está autenticado. Execute: gh auth login — %w", err)
		}
		return fmt.Errorf("gh repo create: %w (%s)", err, msg)
	}
	return nil
}

// Clone clones the given URL into repoDir (parent dir must exist; repoDir must not exist)
func (rs *RepoService) Clone(repoURL string) error {
	parent := filepath.Dir(rs.repoDir)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}
	cmd := exec.Command("git", "clone", repoURL, rs.repoDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Pull runs git pull in repoDir
func (rs *RepoService) Pull() error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = rs.repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull: %w (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// RestoreToHome copies paths from repoDir to homeDir (overwrites). Call BackupExistingConfig before this if needed.
func (rs *RepoService) RestoreToHome(homeDir string, paths []string) error {
	home, err := expandHome(homeDir)
	if err != nil {
		return err
	}
	for _, p := range paths {
		src := filepath.Join(rs.repoDir, p)
		dst := filepath.Join(home, p)
		if err := copyPath(src, dst); err != nil {
			return fmt.Errorf("restore %s: %w", p, err)
		}
	}
	return nil
}

func expandHome(path string) (string, error) {
	if path == "" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return filepath.Abs(path)
}

// copyPath copies a file or directory from src to dst (recursive for dirs)
func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
