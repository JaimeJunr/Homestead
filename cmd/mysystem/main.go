package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jaime/mysystem/internal/app/services"
	"github.com/jaime/mysystem/internal/infrastructure/executor"
	"github.com/jaime/mysystem/internal/infrastructure/installer"
	"github.com/jaime/mysystem/internal/infrastructure/repository"
	"github.com/jaime/mysystem/internal/tui"
)

func main() {
	// Dependency Injection (Manual Wiring)

	// Infrastructure layer - Scripts
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExecutor := executor.NewBashExecutor()

	// Infrastructure layer - Packages/Installers
	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()

	// Application layer
	scriptService := services.NewScriptService(scriptRepo, scriptExecutor)
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	// Presentation layer
	model := tui.NewModel(scriptService, installerService)

	// Create the TUI program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar MySystem: %v\n", err)
		os.Exit(1)
	}
}
