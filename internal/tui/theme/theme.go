// Package theme holds shared Lipgloss styles for the TUI.
package theme

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func StripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	ConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

	Yes = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	No = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	Selected = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 1)

	ScriptScreenOuter = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2)

	ScriptScreenAccent = lipgloss.NewStyle().
				Foreground(lipgloss.Color("63")).
				Bold(true)

	ScriptLogArea = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252"))

	ScriptScreenFooterBar = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Background(lipgloss.Color("235")).
				Padding(0, 1)
)

func InstallerBreadcrumb(segment string) string {
	return "📦 Instaladores > " + segment
}

func InstallerPackageSectionTitle(c types.PackageCategory) string {
	switch c {
	case types.PackageCategoryIDE:
		return "💻 IDEs e Editores"
	case types.PackageCategoryTool:
		return "🔧 Ferramentas de Desenvolvimento"
	case types.PackageCategoryApp:
		return "📱 Aplicações"
	case types.PackageCategoryZshCore:
		return "🐚 Componentes Core (Zsh)"
	case types.PackageCategoryTerminal:
		return "🖥️ Emuladores de Terminal"
	case types.PackageCategoryShell:
		return "🐚 Shells Alternativos"
	case types.PackageCategoryAI:
		return "🤖 Integração com IA"
	case types.PackageCategoryGames:
		return "🎮 Games"
	case types.PackageCategorySysAdmin:
		return "🛡️ Administração de sistemas"
	default:
		return "📦"
	}
}
