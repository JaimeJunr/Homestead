package types

// ConfigScope represents the scope of a shell configuration
type ConfigScope string

const (
	ConfigScopeGeneral ConfigScope = "general"
	ConfigScopeProject ConfigScope = "project"
	ConfigScopeTool    ConfigScope = "tool"
)

// IsValid checks if the config scope is valid
func (c ConfigScope) IsValid() bool {
	switch c {
	case ConfigScopeGeneral, ConfigScopeProject, ConfigScopeTool:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (c ConfigScope) String() string {
	return string(c)
}
