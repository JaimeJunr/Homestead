// Package msg holds Bubble Tea message types for the Homestead TUI.
package msg

import (
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
)

// Progress wraps installer progress updates.
type Progress interfaces.InstallProgress

// InstallComplete signals installation flow finished.
type InstallComplete struct {
	Err error
}

// ZshCoreInstalled is sent when oh-my-zsh detection finishes.
type ZshCoreInstalled struct {
	Installed bool
}

// ZshApplyResult is sent when ApplyConfig finishes.
type ZshApplyResult struct {
	Err error
}

// ZshApplyReturnToMenu is sent after a delay to return to main menu.
type ZshApplyReturnToMenu struct{}

// ScriptCaptured carries stdout/stderr after ExecuteScriptCapture.
type ScriptCaptured struct {
	Output string
	Err    error
}

// ScriptExecFinished is sent after tea.ExecProcess (sudo scripts).
type ScriptExecFinished struct {
	Err error
}

// URLActionDone reports open/copy URL results (Verb: "open" | "copy").
type URLActionDone struct {
	Err  error
	Verb string
}

// ClearKeyboardToast clears the transient keyboard hint line.
type ClearKeyboardToast struct{}

// CatalogFetched is sent after a background fetch of the remote installer catalog.
type CatalogFetched struct {
	Err error
	Ok  bool
}

// NativeMonitorReload carries a snapshot refresh for integrated monitors.
type NativeMonitorReload struct {
	Kind    string
	Battery *monitoring.BatterySnapshot
	Memory  *monitoring.MemorySnapshot
	Err     error
}

// NativeMonitorTick schedules the next monitor refresh.
type NativeMonitorTick struct{}
