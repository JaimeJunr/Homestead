package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepoService(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewRepoService(dir)
	if err != nil {
		t.Fatalf("NewRepoService() error = %v", err)
	}
	if svc.RepoDir() != dir {
		t.Errorf("RepoDir() = %s, want %s", svc.RepoDir(), dir)
	}
}

func TestRepoService_InitRepo_IsRepo(t *testing.T) {
	dir := t.TempDir()
	svc, err := NewRepoService(dir)
	if err != nil {
		t.Fatalf("NewRepoService() error = %v", err)
	}
	if svc.IsRepo() {
		t.Error("IsRepo() = true before InitRepo, want false")
	}
	if err := svc.InitRepo(); err != nil {
		t.Fatalf("InitRepo() error = %v", err)
	}
	if !svc.IsRepo() {
		t.Error("IsRepo() = false after InitRepo, want true")
	}
}

func TestRepoService_CopyToRepo_RestoreToHome(t *testing.T) {
	homeDir := t.TempDir()
	repoDir := t.TempDir()

	// Create source files in "home"
	zshrc := filepath.Join(homeDir, ".zshrc")
	if err := os.WriteFile(zshrc, []byte("# test zshrc\n"), 0644); err != nil {
		t.Fatalf("write .zshrc: %v", err)
	}
	zshDir := filepath.Join(homeDir, ".zsh", "general")
	if err := os.MkdirAll(zshDir, 0755); err != nil {
		t.Fatalf("mkdir .zsh: %v", err)
	}
	if err := os.WriteFile(filepath.Join(zshDir, "aliases.zsh"), []byte("alias ll='ls -la'\n"), 0644); err != nil {
		t.Fatalf("write aliases: %v", err)
	}

	svc, err := NewRepoService(repoDir)
	if err != nil {
		t.Fatalf("NewRepoService() error = %v", err)
	}
	if err := svc.CopyToRepo(homeDir, []string{".zshrc", ".zsh"}); err != nil {
		t.Fatalf("CopyToRepo() error = %v", err)
	}

	// Verify repo has files
	if _, err := os.Stat(filepath.Join(repoDir, ".zshrc")); err != nil {
		t.Errorf("repo .zshrc missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repoDir, ".zsh", "general", "aliases.zsh")); err != nil {
		t.Errorf("repo .zsh/general/aliases.zsh missing: %v", err)
	}

	// Restore to a new "home"
	restoreHome := t.TempDir()
	if err := svc.RestoreToHome(restoreHome, []string{".zshrc", ".zsh"}); err != nil {
		t.Fatalf("RestoreToHome() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(restoreHome, ".zshrc"))
	if err != nil {
		t.Fatalf("read restored .zshrc: %v", err)
	}
	if string(data) != "# test zshrc\n" {
		t.Errorf("restored .zshrc = %q, want # test zshrc\n", string(data))
	}
	data2, _ := os.ReadFile(filepath.Join(restoreHome, ".zsh", "general", "aliases.zsh"))
	if string(data2) != "alias ll='ls -la'\n" {
		t.Errorf("restored aliases.zsh = %q", string(data2))
	}
}

func TestRepoService_ExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("UserHomeDir: %v", err)
	}
	svc, err := NewRepoService("~/homestead-dotfiles")
	if err != nil {
		t.Fatalf("NewRepoService(~/...) error = %v", err)
	}
	want := filepath.Join(home, "homestead-dotfiles")
	if svc.RepoDir() != want {
		t.Errorf("RepoDir() = %s, want %s", svc.RepoDir(), want)
	}
}
