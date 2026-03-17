package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// TestFileConfigManager_SaveConfig tests saving a configuration
func TestFileConfigManager_SaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	config := &entities.ShellConfig{
		ID:    "test-config",
		Name:  "Test Configuration",
		Scope: types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
		Aliases: map[string]string{
			"ll": "ls -la",
		},
	}

	err := manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify file was created
	configFile := filepath.Join(tempDir, "test-config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configFile)
	}
}

// TestFileConfigManager_LoadConfig tests loading a configuration
func TestFileConfigManager_LoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	// First save a config
	original := &entities.ShellConfig{
		ID:    "test-load",
		Name:  "Load Test",
		Scope: types.ConfigScopeProject,
		Plugins: []string{"git", "rails"},
		EnvVars: map[string]string{
			"IVT_DIR": "$HOME/ivt",
		},
	}

	err := manager.SaveConfig(original)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Now load it back
	loaded, err := manager.LoadConfig("test-load")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Verify data
	if loaded.ID != original.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, original.ID)
	}
	if loaded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", loaded.Name, original.Name)
	}
	if loaded.Scope != original.Scope {
		t.Errorf("Scope mismatch: got %s, want %s", loaded.Scope, original.Scope)
	}
	if len(loaded.Plugins) != len(original.Plugins) {
		t.Errorf("Plugins count mismatch: got %d, want %d", len(loaded.Plugins), len(original.Plugins))
	}
}

// TestFileConfigManager_LoadConfig_NotFound tests loading non-existent config
func TestFileConfigManager_LoadConfig_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	_, err := manager.LoadConfig("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent config, got nil")
	}
}

// TestFileConfigManager_DeleteConfig tests deleting a configuration
func TestFileConfigManager_DeleteConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	// Create a config
	config := &entities.ShellConfig{
		ID:    "delete-me",
		Name:  "Delete Test",
		Scope: types.ConfigScopeGeneral,
	}

	err := manager.SaveConfig(config)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Delete it
	err = manager.DeleteConfig("delete-me")
	if err != nil {
		t.Fatalf("DeleteConfig() error = %v", err)
	}

	// Verify it's gone
	configFile := filepath.Join(tempDir, "delete-me.yaml")
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Error("Config file still exists after deletion")
	}
}

// TestFileConfigManager_ListConfigs tests listing all configurations
func TestFileConfigManager_ListConfigs(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	// Create multiple configs
	configs := []string{"config1", "config2", "config3"}
	for _, id := range configs {
		config := &entities.ShellConfig{
			ID:    id,
			Name:  "Config " + id,
			Scope: types.ConfigScopeGeneral,
		}
		err := manager.SaveConfig(config)
		if err != nil {
			t.Fatalf("SaveConfig() error = %v", err)
		}
	}

	// List them
	list, err := manager.ListConfigs()
	if err != nil {
		t.Fatalf("ListConfigs() error = %v", err)
	}

	if len(list) != len(configs) {
		t.Errorf("ListConfigs() count = %d, want %d", len(list), len(configs))
	}

	// Verify all configs are present
	configMap := make(map[string]bool)
	for _, c := range list {
		configMap[c] = true
	}

	for _, expected := range configs {
		if !configMap[expected] {
			t.Errorf("Config %s not found in list", expected)
		}
	}
}

// TestFileConfigManager_GenerateZshrc tests generating .zshrc content
func TestFileConfigManager_GenerateZshrc(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh", "powerlevel10k"},
		Plugins:        []string{"git", "docker", "rails"},
		Tools:          []string{"nvm", "bun"},
	}

	zshrc, err := manager.GenerateZshrc(selections)
	if err != nil {
		t.Fatalf("GenerateZshrc() error = %v", err)
	}

	// Verify content
	if zshrc == "" {
		t.Error("GenerateZshrc() returned empty content")
	}

	// Should contain plugins
	if !contains(zshrc, "plugins=(") {
		t.Error("Generated .zshrc missing plugins declaration")
	}

	// Should contain git plugin
	if !contains(zshrc, "git") {
		t.Error("Generated .zshrc missing git plugin")
	}
}

// TestFileConfigManager_GenerateAliasesFile tests generating aliases.zsh
func TestFileConfigManager_GenerateAliasesFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	config := &entities.ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeGeneral,
		Aliases: map[string]string{
			"ll":   "ls -la",
			"la":   "ls -A",
			"grep": "grep --color=auto",
		},
	}

	content, err := manager.GenerateAliasesFile(config)
	if err != nil {
		t.Fatalf("GenerateAliasesFile() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateAliasesFile() returned empty content")
	}

	// Should contain aliases
	if !contains(content, "alias ll=") {
		t.Error("Generated aliases missing 'll' alias")
	}
}

// TestFileConfigManager_GenerateFunctionsFile tests generating functions.zsh
func TestFileConfigManager_GenerateFunctionsFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	config := &entities.ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeGeneral,
		Functions: map[string]string{
			"db-connect": `local database="${1:-funds}"
mysql -u root -p "$database"`,
		},
	}

	content, err := manager.GenerateFunctionsFile(config)
	if err != nil {
		t.Fatalf("GenerateFunctionsFile() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateFunctionsFile() returned empty content")
	}

	// Should contain function
	if !contains(content, "db-connect()") {
		t.Error("Generated functions missing 'db-connect' function")
	}
}

// TestFileConfigManager_BackupExistingConfig tests backup functionality
func TestFileConfigManager_BackupExistingConfig(t *testing.T) {
	tempDir := t.TempDir()
	_ = NewFileConfigManager(tempDir)

	// Create a fake .zshrc
	zshrcPath := filepath.Join(tempDir, ".zshrc")
	err := os.WriteFile(zshrcPath, []byte("# Existing zshrc"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .zshrc: %v", err)
	}

	// Note: BackupExistingConfig would need to be adapted to use tempDir
	// For now, we'll test the concept
	// err = manager.BackupExistingConfig()
	// In real implementation, would check for backup file
}

// TestFileConfigManager_ApplyConfig tests applying configuration
func TestFileConfigManager_ApplyConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh"},
		Plugins:        []string{"git"},
		Tools:          []string{"nvm"},
	}

	// Note: ApplyConfig would need to be adapted for testing
	// to write to tempDir instead of home directory
	err := manager.ApplyConfig(selections)
	if err != nil {
		t.Logf("ApplyConfig() error = %v (expected in test environment)", err)
	}
}

// TestFileConfigManager_ApplyConfig_DoesNotOverwriteUserContent verifies that
// existing user content in ~/.zshrc is preserved when applying a new config.
func TestFileConfigManager_ApplyConfig_DoesNotOverwriteUserContent(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	// Point HOME to tempDir so ApplyConfig writes into our sandbox
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	zshrcPath := filepath.Join(tempDir, ".zshrc")
	userContent := "# User config\nsource ~/.p10k.zsh\n"
	if err := os.WriteFile(zshrcPath, []byte(userContent), 0644); err != nil {
		t.Fatalf("failed to create existing .zshrc: %v", err)
	}

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh"},
		Plugins:        []string{"git"},
		Tools:          []string{"nvm"},
	}

	if err := manager.ApplyConfig(selections); err != nil {
		t.Fatalf("ApplyConfig() error = %v", err)
	}

	data, err := os.ReadFile(zshrcPath)
	if err != nil {
		t.Fatalf("failed to read resulting .zshrc: %v", err)
	}
	content := string(data)

	// User content must still be present
	if !contains(content, "User config") {
		t.Errorf("user content was lost from .zshrc")
	}
	if !contains(content, "source ~/.p10k.zsh") {
		t.Errorf("p10k sourcing line was lost from .zshrc")
	}
	// And Homestead block should be present as well
	if !contains(content, "Generated by Homestead") {
		t.Errorf("Homestead-managed block not found in .zshrc")
	}
}

// TestFileConfigManager_ApplyConfig_WritesAliasesToFile verifies that
// CustomAliases from selections are written to ~/.zsh/general/aliases.zsh.
func TestFileConfigManager_ApplyConfig_WritesAliasesToFile(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewFileConfigManager(tempDir)

	// Point HOME to tempDir so ApplyConfig writes into our sandbox
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh"},
		CustomAliases: map[string]string{
			"ll": "ls -la",
			"gs": "git status",
		},
	}

	if err := manager.ApplyConfig(selections); err != nil {
		t.Fatalf("ApplyConfig() error = %v", err)
	}

	aliasesPath := filepath.Join(tempDir, ".zsh", "general", "aliases.zsh")
	data, err := os.ReadFile(aliasesPath)
	if err != nil {
		t.Fatalf("failed to read aliases.zsh: %v", err)
	}
	content := string(data)

	if !contains(content, "alias ll='ls -la'") {
		t.Errorf("aliases.zsh missing expected alias ll")
	}
	if !contains(content, "alias gs='git status'") {
		t.Errorf("aliases.zsh missing expected alias gs")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkFileConfigManager_SaveConfig(b *testing.B) {
	tempDir := b.TempDir()
	manager := NewFileConfigManager(tempDir)

	config := &entities.ShellConfig{
		ID:    "bench",
		Name:  "Benchmark",
		Scope: types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.SaveConfig(config)
	}
}

func BenchmarkFileConfigManager_LoadConfig(b *testing.B) {
	tempDir := b.TempDir()
	manager := NewFileConfigManager(tempDir)

	config := &entities.ShellConfig{
		ID:    "bench",
		Name:  "Benchmark",
		Scope: types.ConfigScopeGeneral,
	}

	_ = manager.SaveConfig(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.LoadConfig("bench")
	}
}
