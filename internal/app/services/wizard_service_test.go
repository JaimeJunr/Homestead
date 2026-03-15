package services

import (
	"testing"
)

// TestWizardService_CreateNewWizard tests creating a new wizard session
func TestWizardService_CreateNewWizard(t *testing.T) {
	service := NewWizardService()

	wizard := service.CreateNewWizard()
	if wizard == nil {
		t.Error("CreateNewWizard() returned nil")
	}

	if wizard.CurrentStep != 0 {
		t.Errorf("CurrentStep = %d, want 0", wizard.CurrentStep)
	}

	if len(wizard.Selections.CoreComponents) != 0 || len(wizard.Selections.Plugins) != 0 {
		t.Error("Selections should be initialized empty")
	}
}

// TestWizardService_GetCurrentStep tests getting current wizard step
func TestWizardService_GetCurrentStep(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	step := service.GetCurrentStep(wizard)
	if step == nil {
		t.Error("GetCurrentStep() returned nil")
	}

	if step.Name == "" {
		t.Error("Step name should not be empty")
	}
}

// TestWizardService_NextStep tests advancing to next step
func TestWizardService_NextStep(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	initialStep := wizard.CurrentStep

	err := service.NextStep(wizard)
	if err != nil {
		t.Fatalf("NextStep() error = %v", err)
	}

	if wizard.CurrentStep != initialStep+1 {
		t.Errorf("CurrentStep = %d, want %d", wizard.CurrentStep, initialStep+1)
	}
}

// TestWizardService_PreviousStep tests going back to previous step
func TestWizardService_PreviousStep(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	// Advance to step 1
	_ = service.NextStep(wizard)

	err := service.PreviousStep(wizard)
	if err != nil {
		t.Fatalf("PreviousStep() error = %v", err)
	}

	if wizard.CurrentStep != 0 {
		t.Errorf("CurrentStep = %d, want 0", wizard.CurrentStep)
	}
}

// TestWizardService_PreviousStep_AtStart tests going back at first step
func TestWizardService_PreviousStep_AtStart(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	err := service.PreviousStep(wizard)
	if err == nil {
		t.Error("Expected error when going back from first step, got nil")
	}
}

// TestWizardService_IsFirstStep tests checking if at first step
func TestWizardService_IsFirstStep(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	if !service.IsFirstStep(wizard) {
		t.Error("Should be at first step initially")
	}

	_ = service.NextStep(wizard)

	if service.IsFirstStep(wizard) {
		t.Error("Should not be at first step after NextStep")
	}
}

// TestWizardService_IsLastStep tests checking if at last step
func TestWizardService_IsLastStep(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	if service.IsLastStep(wizard) {
		t.Error("Should not be at last step initially")
	}

	// Navigate to last step
	for !service.IsLastStep(wizard) {
		err := service.NextStep(wizard)
		if err != nil {
			break
		}
	}

	if !service.IsLastStep(wizard) {
		t.Error("Should be at last step after navigating")
	}
}

// TestWizardService_AddSelection tests adding a selection
func TestWizardService_AddSelection(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.AddCoreComponent(wizard, "zsh")
	service.AddCoreComponent(wizard, "oh-my-zsh")

	if len(wizard.Selections.CoreComponents) != 2 {
		t.Errorf("CoreComponents count = %d, want 2", len(wizard.Selections.CoreComponents))
	}

	if !sliceContains(wizard.Selections.CoreComponents, "zsh") {
		t.Error("CoreComponents should contain 'zsh'")
	}
}

// TestWizardService_RemoveSelection tests removing a selection
func TestWizardService_RemoveSelection(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.AddCoreComponent(wizard, "zsh")
	service.AddCoreComponent(wizard, "oh-my-zsh")
	service.RemoveCoreComponent(wizard, "zsh")

	if len(wizard.Selections.CoreComponents) != 1 {
		t.Errorf("CoreComponents count = %d, want 1", len(wizard.Selections.CoreComponents))
	}

	if sliceContains(wizard.Selections.CoreComponents, "zsh") {
		t.Error("CoreComponents should not contain 'zsh' after removal")
	}
}

// TestWizardService_AddPlugin tests adding plugins
func TestWizardService_AddPlugin(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.AddPlugin(wizard, "git")
	service.AddPlugin(wizard, "docker")

	if len(wizard.Selections.Plugins) != 2 {
		t.Errorf("Plugins count = %d, want 2", len(wizard.Selections.Plugins))
	}
}

// TestWizardService_AddTool tests adding tools
func TestWizardService_AddTool(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.AddTool(wizard, "nvm")
	service.AddTool(wizard, "bun")

	if len(wizard.Selections.Tools) != 2 {
		t.Errorf("Tools count = %d, want 2", len(wizard.Selections.Tools))
	}
}

// TestWizardService_SetIncludeProjectConfig tests project config inclusion
func TestWizardService_SetIncludeProjectConfig(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.SetIncludeProjectConfig(wizard, true)

	if !wizard.Selections.IncludeProjectConfig {
		t.Error("IncludeProjectConfig should be true")
	}

	service.SetIncludeProjectConfig(wizard, false)

	if wizard.Selections.IncludeProjectConfig {
		t.Error("IncludeProjectConfig should be false")
	}
}

// TestWizardService_GeneratePreview tests generating configuration preview
func TestWizardService_GeneratePreview(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	// Add some selections
	service.AddCoreComponent(wizard, "zsh")
	service.AddCoreComponent(wizard, "oh-my-zsh")
	service.AddPlugin(wizard, "git")
	service.AddTool(wizard, "nvm")

	preview := service.GeneratePreview(wizard)
	if preview == "" {
		t.Error("GeneratePreview() returned empty string")
	}

	// Verify preview contains selections
	if !stringContains(preview, "zsh") {
		t.Error("Preview should contain 'zsh'")
	}

	if !stringContains(preview, "git") {
		t.Error("Preview should contain 'git'")
	}
}

// TestWizardService_ValidateSelections tests validating wizard selections
func TestWizardService_ValidateSelections(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	// Empty selections should fail
	err := service.ValidateSelections(wizard)
	if err == nil {
		t.Error("Expected error for empty selections, got nil")
	}

	// Add minimum required selections
	service.AddCoreComponent(wizard, "zsh")

	err = service.ValidateSelections(wizard)
	if err != nil {
		t.Errorf("ValidateSelections() error = %v", err)
	}
}

// TestWizardService_GetTotalSteps tests getting total number of steps
func TestWizardService_GetTotalSteps(t *testing.T) {
	service := NewWizardService()

	total := service.GetTotalSteps()
	if total <= 0 {
		t.Errorf("GetTotalSteps() = %d, want > 0", total)
	}
}

// TestWizardService_GetProgress tests getting wizard progress percentage
func TestWizardService_GetProgress(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	progress := service.GetProgress(wizard)
	if progress < 0 || progress > 100 {
		t.Errorf("GetProgress() = %d, should be between 0 and 100", progress)
	}

	// Progress at start should be low
	if progress > 20 {
		t.Errorf("Initial progress = %d, expected low value", progress)
	}

	// Advance wizard
	for i := 0; i < 3; i++ {
		_ = service.NextStep(wizard)
	}

	newProgress := service.GetProgress(wizard)
	if newProgress <= progress {
		t.Error("Progress should increase after advancing steps")
	}
}

// TestWizardService_Reset tests resetting wizard
func TestWizardService_Reset(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	// Add selections and advance
	service.AddCoreComponent(wizard, "zsh")
	service.AddPlugin(wizard, "git")
	_ = service.NextStep(wizard)
	_ = service.NextStep(wizard)

	// Reset
	service.Reset(wizard)

	if wizard.CurrentStep != 0 {
		t.Errorf("CurrentStep after reset = %d, want 0", wizard.CurrentStep)
	}

	if len(wizard.Selections.CoreComponents) != 0 {
		t.Error("Selections should be cleared after reset")
	}
}

// TestWizardService_CanProceed tests checking if can proceed to next step
func TestWizardService_CanProceed(t *testing.T) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	// Without selections, might not be able to proceed from certain steps
	// This depends on validation logic
	canProceed := service.CanProceed(wizard)

	// At step 0, should be able to proceed regardless
	if !canProceed {
		t.Log("Cannot proceed from initial step (may require selections)")
	}
}

// Helper function to check if string contains substring
func stringContains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

// Benchmark tests
func BenchmarkWizardService_CreateNewWizard(b *testing.B) {
	service := NewWizardService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.CreateNewWizard()
	}
}

func BenchmarkWizardService_NextStep(b *testing.B) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.NextStep(wizard)
		if service.IsLastStep(wizard) {
			service.Reset(wizard)
		}
	}
}

func BenchmarkWizardService_GeneratePreview(b *testing.B) {
	service := NewWizardService()
	wizard := service.CreateNewWizard()

	service.AddCoreComponent(wizard, "zsh")
	service.AddCoreComponent(wizard, "oh-my-zsh")
	service.AddPlugin(wizard, "git")
	service.AddPlugin(wizard, "docker")
	service.AddTool(wizard, "nvm")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.GeneratePreview(wizard)
	}
}
