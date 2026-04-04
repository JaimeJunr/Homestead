package interfaces

import (
	"os/exec"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
)

// ScriptExecutor defines the interface for script execution
type ScriptExecutor interface {
	// Execute executes a script
	Execute(script *entities.Script) error

	// ExecuteCapture runs the script with stdout/stderr captured (no TTY).
	ExecuteCapture(script *entities.Script) (output string, err error)

	// InteractiveCommand returns a cmd for tea.ExecProcess (sudo/password, full terminal).
	InteractiveCommand(script *entities.Script) (*exec.Cmd, error)

	// CanExecute checks if a script can be executed
	CanExecute(script *entities.Script) bool

	// Validate validates a script before execution
	Validate(script *entities.Script) error
}
