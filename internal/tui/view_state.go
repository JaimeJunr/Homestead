package tui

// ViewState represents different views in the TUI
type ViewState int

const (
	ViewMainMenu ViewState = iota
	ViewScriptList
	ViewInstallerCategories
	ViewPackageList
	ViewConfirmation
	ViewScriptOutput
	ViewNativeMonitor
	ViewInstalling
	ViewZshWizard
	ViewZshApplying
	ViewZshRepoWizard
)

// menuAction identifies the main menu action
const (
	menuActionCleanup    = "cleanup"
	menuActionMonitoring = "monitoring"
	menuActionInstallers = "installers"
	menuActionZshPlugins = "zsh_plugins" // Plugins e temas Zsh (wizard local)
	menuActionZshRepo    = "zsh_repo"    // Configurar Zsh (repo backup/migração)
	menuActionSettings   = "settings"
	menuActionQuit       = "quit"
)
