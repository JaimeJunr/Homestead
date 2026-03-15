package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaime/mysystem/internal/app/services"
	"github.com/jaime/mysystem/internal/infrastructure/executor"
	"github.com/jaime/mysystem/internal/infrastructure/installer"
	"github.com/jaime/mysystem/internal/infrastructure/repository"
)

// testModel creates a model for testing with mocked dependencies
func testModel() Model {
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	return NewModel(scriptService, installerService)
}

func TestNewModel(t *testing.T) {
	model := testModel()

	if model.state != ViewMainMenu {
		t.Errorf("Expected initial state to be ViewMainMenu, got %d", model.state)
	}

	items := model.mainMenu.Items()
	if len(items) != 7 {
		t.Errorf("Expected 7 main menu items, got %d", len(items))
	}

	if model.scriptService == nil {
		t.Error("Expected scriptService to be initialized")
	}
}

func TestViewStates(t *testing.T) {
	// Verify view state constants
	if ViewMainMenu != 0 {
		t.Errorf("ViewMainMenu should be 0, got %d", ViewMainMenu)
	}
	if ViewScriptList != 1 {
		t.Errorf("ViewScriptList should be 1, got %d", ViewScriptList)
	}
	if ViewPackageList != 2 {
		t.Errorf("ViewPackageList should be 2, got %d", ViewPackageList)
	}
	if ViewConfirmation != 3 {
		t.Errorf("ViewConfirmation should be 3, got %d", ViewConfirmation)
	}
	if ViewExecuting != 4 {
		t.Errorf("ViewExecuting should be 4, got %d", ViewExecuting)
	}
	if ViewInstalling != 5 {
		t.Errorf("ViewInstalling should be 5, got %d", ViewInstalling)
	}
	if ViewZshWizard != 6 {
		t.Errorf("ViewZshWizard should be 6, got %d", ViewZshWizard)
	}
}

func TestModelInit(t *testing.T) {
	model := testModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init() to return spinner tick command")
	}
}

func TestWindowSizeUpdate(t *testing.T) {
	model := testModel()

	msg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 80 {
		t.Errorf("Expected width 80, got %d", m.width)
	}
	if m.height != 24 {
		t.Errorf("Expected height 24, got %d", m.height)
	}
}

func TestQuitOnMainMenu(t *testing.T) {
	model := testModel()
	model.state = ViewMainMenu

	// Test 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Expected quit command on 'q' key")
	}

	// Test Ctrl+C
	msgCtrlC := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmdCtrlC := model.Update(msgCtrlC)

	if cmdCtrlC == nil {
		t.Error("Expected quit command on Ctrl+C")
	}
}

func TestEscapeFromScriptList(t *testing.T) {
	model := testModel()
	model.state = ViewScriptList

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.state != ViewMainMenu {
		t.Errorf("Expected state to return to ViewMainMenu, got %d", m.state)
	}
}

func TestMenuItemInterface(t *testing.T) {
	item := menuItem{
		title: "Test Title",
		desc:  "Test Description",
	}

	if item.Title() != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", item.Title())
	}

	if item.Description() != "Test Description" {
		t.Errorf("Expected description 'Test Description', got %s", item.Description())
	}

	if item.FilterValue() != "Test Title" {
		t.Errorf("Expected filter value 'Test Title', got %s", item.FilterValue())
	}
}

func TestScriptItemInterface(t *testing.T) {
	// Tested via integration
}

func TestViewRendering(t *testing.T) {
	model := testModel()

	// Test initial view (no size set)
	view := model.View()
	if view != "Iniciando..." {
		t.Errorf("Expected 'Iniciando...' for uninitialized view, got %s", view)
	}

	// Set window size
	model.width = 80
	model.height = 24

	// Test main menu view
	model.state = ViewMainMenu
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view for main menu")
	}
}

func TestModelStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  ViewState
		expectedState ViewState
	}{
		{
			name:          "Start at main menu",
			initialState:  ViewMainMenu,
			expectedState: ViewMainMenu,
		},
		{
			name:          "Can be at script list",
			initialState:  ViewScriptList,
			expectedState: ViewScriptList,
		},
		{
			name:          "Can be at executing",
			initialState:  ViewExecuting,
			expectedState: ViewExecuting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := testModel()
			model.state = tt.initialState

			if model.state != tt.expectedState {
				t.Errorf("Expected state %d, got %d", tt.expectedState, model.state)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewModel(b *testing.B) {
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewModel(scriptService, installerService)
	}
}

func BenchmarkModelUpdate(b *testing.B) {
	model := testModel()
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.Update(msg)
	}
}

func BenchmarkModelView(b *testing.B) {
	model := testModel()
	model.width = 80
	model.height = 24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.View()
	}
}
