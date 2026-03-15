package types

// PluginSource represents the source/origin of a Zsh plugin
type PluginSource string

const (
	PluginSourceBuiltIn  PluginSource = "builtin"  // Comes with Oh My Zsh
	PluginSourceExternal PluginSource = "external" // External repository (GitHub, etc)
	PluginSourceCustom   PluginSource = "custom"   // User custom plugin
)

// IsValid checks if the plugin source is valid
func (ps PluginSource) IsValid() bool {
	switch ps {
	case PluginSourceBuiltIn, PluginSourceExternal, PluginSourceCustom:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (ps PluginSource) String() string {
	return string(ps)
}
