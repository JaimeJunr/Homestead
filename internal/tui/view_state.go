package tui

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

const (
	menuActionCleanup    = "cleanup"
	menuActionMonitoring = "monitoring"
	menuActionInstallers = "installers"
	menuActionZshPlugins = "zsh_plugins"
	menuActionZshRepo    = "zsh_repo"
	menuActionSettings   = "settings"
	menuActionQuit       = "quit"
)
