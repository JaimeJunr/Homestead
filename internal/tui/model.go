package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
	"github.com/JaimeJunr/Homestead/internal/tui/cmds"
	"github.com/JaimeJunr/Homestead/internal/tui/items"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

// Model is the main TUI model
type Model struct {
	scriptService    *services.ScriptService
	installerService *services.InstallerService
	configService    *services.ConfigService
	repoService      *services.RepoService
	state            ViewState
	mainMenu         list.Model
	scriptList       list.Model
	installerList    list.Model
	packageList      list.Model
	selectedMenu     int
	selectedItem     interface{} // Can be Script or Package
	confirmYes       bool        // true = yes selected, false = no selected
	confirmReturn    ViewState   // tela para voltar se cancelar a confirmação (lista de pacotes/scripts)
	confirmReturnOK  bool        // se false, cancelar volta ao menu principal
	width            int
	height           int
	err              error
	keyboardToast    string // feedback para o/c (abrir/copiar URL) sem mouse

	// Installation progress
	progress       progress.Model
	spinner        spinner.Model
	installStatus  string
	installMessage string
	installPercent float64
	canAbort       bool
	aborted        bool

	// Zsh plugins wizard (Plugins e temas Zsh)
	zshWizard *ZshWizardModel

	// Zsh repo wizard (Configurar Zsh - backup/migração via repositório)
	zshRepoWizard *ZshRepoModel

	// Zsh core: when true, "Plugins e temas Zsh" is shown in menu (oh-my-zsh installed)
	zshCoreInstalled bool
	zshCoreChecked   bool

	// Zsh apply feedback: phase "applying" | "success" | "error"
	zshApplyPhase string
	zshApplyError error

	// Script output (in-TUI); phase "running" | "done"
	scriptOutputView  viewport.Model
	scriptOutputPhase string
	scriptOutputTitle string
	scriptOutputErr   error

	// Monitores integrados (bateria / memória)
	nativeMonitorKind string
	nativeBattery     *monitoring.BatterySnapshot
	nativeBatteryErr  error
	nativeMemory      *monitoring.MemorySnapshot
	nativeMemoryErr   error

	// scriptListParent: para onde Esc volta a partir de ViewScriptList (menu principal ou categorias de instaladores).
	scriptListParent ViewState
	// scriptListAsInstaller: lista de utilitários aberta a partir de Instaladores (UX alinhada a pacotes).
	scriptListAsInstaller bool

	// Installer package list: categories filter for the current ViewPackageList (refresh after remote catalog).
	packageListCategories []types.PackageCategory
	catalogURL            string
}

// NewModel creates the TUI model with dependencies injected.
// catalogURL may be empty to skip remote catalog fetch (e.g. tests).
func NewModel(scriptService *services.ScriptService, installerService *services.InstallerService, configService *services.ConfigService, repoService *services.RepoService, catalogURL string) Model {
	mainItems := getMainMenuItems(false) // will refresh when zsh core check completes
	mainList := list.New(mainItems, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "Homestead - Gerenciador de Sistema"
	mainList.SetShowStatusBar(false)
	mainList.SetFilteringEnabled(false)

	prog := progress.New(progress.WithDefaultGradient())

	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return Model{
		scriptService:         scriptService,
		installerService:      installerService,
		configService:         configService,
		repoService:           repoService,
		catalogURL:            catalogURL,
		state:                 ViewMainMenu,
		mainMenu:              mainList,
		progress:              prog,
		spinner:               spin,
		confirmYes:            false,
		scriptListParent:      ViewMainMenu,
		scriptListAsInstaller: false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	batch := []tea.Cmd{m.spinner.Tick, cmds.CheckZshCoreInstalled(m.installerService)}
	if c := cmds.FetchCatalog(m.catalogURL, m.installerService); c != nil {
		batch = append(batch, c)
	}
	return tea.Batch(batch...)
}

// Update handles messages and updates state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if m.state == ViewScriptOutput && m.scriptOutputPhase == "done" {
			var vcmd tea.Cmd
			m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
			return m, vcmd
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height-4)
		if m.scriptList.Items() != nil {
			m.scriptList.SetSize(msg.Width, msg.Height-4)
		}
		if m.installerList.Items() != nil {
			m.installerList.SetSize(msg.Width, msg.Height-4)
		}
		if m.packageList.Items() != nil {
			m.packageList.SetSize(msg.Width, msg.Height-4)
		}
		if m.state == ViewScriptOutput {
			m.syncScriptOutputViewport()
		}
		return m, nil

	case btmsg.Progress:
		m.installStatus = msg.Status
		m.installMessage = msg.Message
		m.installPercent = float64(msg.Progress) / 100.0
		m.canAbort = msg.CanAbort

		if msg.IsCompleted {
			return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return btmsg.InstallComplete{Err: msg.Error}
			})
		}
		return m, nil

	case btmsg.InstallComplete:
		m.state = ViewMainMenu
		m.aborted = false
		return m, cmds.CheckZshCoreInstalled(m.installerService)

	case btmsg.ZshCoreInstalled:
		m.zshCoreChecked = true
		m.zshCoreInstalled = msg.Installed
		m.mainMenu.SetItems(getMainMenuItems(m.zshCoreInstalled))
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.state == ViewScriptOutput {
			if m.scriptOutputPhase == "done" {
				switch msg.String() {
				case "enter", "esc", "q":
					m.state = m.confirmReturn
					m.scriptOutputPhase = ""
					m.scriptOutputTitle = ""
					m.scriptOutputErr = nil
					return m, nil
				}
				var vcmd tea.Cmd
				m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
				return m, vcmd
			}
			return m, nil
		}
		if m.state == ViewNativeMonitor {
			switch msg.String() {
			case "enter", "esc", "q":
				m.state = m.confirmReturn
				m.nativeMonitorKind = ""
				m.nativeBattery, m.nativeMemory = nil, nil
				m.nativeBatteryErr, m.nativeMemoryErr = nil, nil
				return m, nil
			case "r":
				return m, m.nativeMonitorLoadCmd()
			}
			return m, nil
		}
		if m.state == ViewScriptList && m.err != nil {
			m.err = nil
		}
		if m.state == ViewZshApplying && (m.zshApplyPhase == "success" || m.zshApplyPhase == "error") {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.state = ViewMainMenu
				m.zshApplyPhase = ""
				m.zshApplyError = nil
				return m, nil
			}
		}
		if m.state != ViewZshWizard && m.state != ViewZshRepoWizard {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.state == ViewMainMenu {
					return m, tea.Quit
				}
				if m.state == ViewInstalling && m.canAbort {
					m.aborted = true
					m.installMessage = "Instalação abortada pelo usuário"
					m.state = ViewMainMenu
					return m, nil
				}
			case "esc":
				switch m.state {
				case ViewScriptList:
					m.state = m.scriptListParent
					m.confirmYes = false
					return m, nil
				case ViewConfirmation:
					if m.confirmReturnOK {
						m.state = m.confirmReturn
					} else {
						m.state = ViewMainMenu
					}
					m.confirmReturnOK = false
					m.confirmYes = false
					return m, nil
				case ViewPackageList:
					m.state = ViewInstallerCategories
					return m, nil
				case ViewInstallerCategories:
					m.state = ViewMainMenu
					return m, nil
				}
			case "left", "h":
				if m.state == ViewConfirmation {
					m.confirmYes = false
					return m, nil
				}
			case "right", "l":
				if m.state == ViewConfirmation {
					m.confirmYes = true
					return m, nil
				}
			case "o", "O":
				return m.handleURLShortcut(false)
			case "c", "C":
				return m.handleURLShortcut(true)
			case "enter":
				return m.handleEnter()
			}
		}

	case btmsg.ZshApplyResult:
		if m.state == ViewZshApplying {
			if msg.Err != nil {
				m.zshApplyPhase = "error"
				m.zshApplyError = msg.Err
			} else {
				m.zshApplyPhase = "success"
				m.zshApplyError = nil
			}
			return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
				return btmsg.ZshApplyReturnToMenu{}
			})
		}
		return m, nil

	case btmsg.ZshApplyReturnToMenu:
		if m.state == ViewZshApplying {
			m.state = ViewMainMenu
			m.zshApplyPhase = ""
			m.zshApplyError = nil
		}
		return m, nil

	case btmsg.ScriptCaptured:
		if m.state != ViewScriptOutput {
			return m, nil
		}
		m.scriptOutputPhase = "done"
		m.scriptOutputErr = msg.Err
		text := theme.StripANSI(msg.Output)
		if strings.TrimSpace(text) == "" {
			text = "(sem saída no stdout/stderr)"
		}
		if msg.Err != nil {
			text += "\n\n──\n" + msg.Err.Error()
		}
		m.scriptOutputView.SetContent(text)
		m.scriptOutputView.GotoTop()
		return m, nil

	case btmsg.ScriptExecFinished:
		m.state = m.confirmReturn
		m.scriptOutputPhase = ""
		m.scriptOutputTitle = ""
		m.scriptOutputErr = nil
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
		}
		return m, nil

	case btmsg.NativeMonitorReload:
		if m.state != ViewNativeMonitor || msg.Kind != m.nativeMonitorKind {
			return m, nil
		}
		switch msg.Kind {
		case entities.NativeMonitorBattery:
			m.nativeBattery = msg.Battery
			m.nativeBatteryErr = msg.Err
		case entities.NativeMonitorMemory:
			m.nativeMemory = msg.Memory
			m.nativeMemoryErr = msg.Err
		}
		return m, nativeMonitorScheduleTick()

	case btmsg.NativeMonitorTick:
		if m.state != ViewNativeMonitor {
			return m, nil
		}
		return m, m.nativeMonitorLoadCmd()

	case btmsg.CatalogFetched:
		if msg.Err != nil {
			return m, nil
		}
		var nextCmd tea.Cmd
		if msg.Ok {
			if m.state == ViewPackageList && len(m.packageListCategories) > 0 {
				sel := m.packageList.Index()
				m.loadPackagesFromCategories(m.packageListCategories)
				pkgRows := m.packageList.Items()
				if len(pkgRows) > 0 {
					if sel < 0 {
						sel = 0
					}
					if sel >= len(pkgRows) {
						sel = len(pkgRows) - 1
					}
					m.packageList.Select(sel)
				}
			}
			m.keyboardToast = "Catálogo de instaladores atualizado."
			nextCmd = tea.Tick(2*time.Second, func(time.Time) tea.Msg { return btmsg.ClearKeyboardToast{} })
		}
		return m, nextCmd

	case btmsg.URLActionDone:
		if msg.Err != nil {
			m.keyboardToast = fmt.Sprintf("⚠ %v", msg.Err)
		} else if msg.Verb == "copy" {
			m.keyboardToast = "URL copiada para a área de transferência."
		} else {
			m.keyboardToast = "URL aberta no navegador (app padrão)."
		}
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return btmsg.ClearKeyboardToast{} })

	case btmsg.ClearKeyboardToast:
		m.keyboardToast = ""
		return m, nil
	}

	if m.state == ViewZshWizard && m.zshWizard != nil {
		newWizard, cmd := m.zshWizard.Update(msg)
		wizard := newWizard.(ZshWizardModel)
		m.zshWizard = &wizard

		if wizard.IsDone() || wizard.IsCancelled() {
			if wizard.IsCancelled() {
				m.state = ViewMainMenu
				m.zshWizard = nil
				return m, cmd
			}
			selections := wizard.GetSelections()
			m.zshWizard = nil
			m.state = ViewZshApplying
			m.zshApplyPhase = "applying"
			m.zshApplyError = nil
			return m, cmds.ApplyZshConfig(m.configService, selections)
		}

		return m, cmd
	}

	if m.state == ViewZshRepoWizard && m.zshRepoWizard != nil {
		newRepo, cmd := m.zshRepoWizard.Update(msg)
		repoWizard := newRepo.(ZshRepoModel)
		m.zshRepoWizard = &repoWizard

		if repoWizard.IsDone() || repoWizard.IsCancelled() {
			m.state = ViewMainMenu
			m.zshRepoWizard = nil
			return m, cmd
		}
		return m, cmd
	}

	var cmd tea.Cmd
	switch m.state {
	case ViewMainMenu:
		m.mainMenu, cmd = m.mainMenu.Update(msg)
	case ViewScriptList:
		m.scriptList, cmd = m.scriptList.Update(msg)
	case ViewInstallerCategories:
		m.installerList, cmd = m.installerList.Update(msg)
	case ViewPackageList:
		m.packageList, cmd = m.packageList.Update(msg)
	}

	return m, cmd
}

// handleEnter handles the enter key based on current state
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewMainMenu:
		selected := m.mainMenu.SelectedItem()
		item, ok := selected.(items.MenuItem)
		if !ok {
			return m, nil
		}
		switch item.Action {
		case menuActionCleanup:
			m.state = ViewScriptList
			m.selectedMenu = 0
			m.loadScripts(types.CategoryCleanup)
		case menuActionMonitoring:
			m.state = ViewScriptList
			m.selectedMenu = 1
			m.loadScripts(types.CategoryMonitoring)
		case menuActionInstallers:
			m.state = ViewInstallerCategories
			m.selectedMenu = 2
			m.loadInstallerCategories()
		case menuActionZshPlugins:
			m.state = ViewZshWizard
			wizardService := services.NewWizardService()
			wizard := NewZshWizardModel(wizardService)
			wizard.width = m.width
			wizard.height = m.height
			m.zshWizard = &wizard
		case menuActionZshRepo:
			m.state = ViewZshRepoWizard
			repoWizard := NewZshRepoModel(m.repoService, m.configService)
			repoWizard.width = m.width
			repoWizard.height = m.height
			m.zshRepoWizard = &repoWizard
		case menuActionQuit:
			return m, tea.Quit
		case menuActionSettings:
			return m, nil
		default:
			return m, nil
		}

	case ViewScriptList:
		selected := m.scriptList.SelectedItem()
		if scriptItem, ok := selected.(items.ScriptItem); ok {
			m.selectedItem = scriptItem.Script
			m.state = ViewConfirmation
			m.confirmYes = false
			m.confirmReturn = ViewScriptList
			m.confirmReturnOK = true
		}

	case ViewPackageList:
		selected := m.packageList.SelectedItem()
		if pkgItem, ok := selected.(items.PackageItem); ok {
			m.selectedItem = pkgItem.Pkg
			m.state = ViewConfirmation
			m.confirmYes = false
			m.confirmReturn = ViewPackageList
			m.confirmReturnOK = true
		}

	case ViewInstallerCategories:
		selected := m.installerList.SelectedItem()
		catItem, ok := selected.(items.InstallerCategoryItem)
		if !ok {
			break
		}
		if catItem.Utilities {
			m.state = ViewScriptList
			m.loadScriptsWithParent(types.CategoryUtilities, ViewInstallerCategories)
			return m, nil
		}
		if len(catItem.Categories) > 0 {
			m.state = ViewPackageList
			m.packageListCategories = append([]types.PackageCategory(nil), catItem.Categories...)
			m.loadPackagesFromCategories(catItem.Categories)
		}

	case ViewConfirmation:
		if m.confirmYes {
			switch item := m.selectedItem.(type) {
			case entities.Script:
				if item.NativeMonitor != "" {
					m.nativeMonitorKind = item.NativeMonitor
					m.nativeBattery, m.nativeMemory = nil, nil
					m.nativeBatteryErr, m.nativeMemoryErr = nil, nil
					m.state = ViewNativeMonitor
					return m, m.nativeMonitorLoadCmd()
				}
				m.scriptOutputTitle = item.Name
				m.scriptOutputPhase = "running"
				m.scriptOutputErr = nil
				m.scriptOutputView = newScriptOutputViewport(m.width, m.height)
				m.state = ViewScriptOutput
				if item.RequiresSudo {
					cmd, err := m.scriptService.ScriptInteractiveCommand(item.ID)
					if err != nil {
						m.state = m.confirmReturn
						m.scriptOutputPhase = ""
						m.scriptOutputTitle = ""
						m.err = err
						return m, nil
					}
					return m, tea.ExecProcess(cmd, func(execErr error) tea.Msg {
						return btmsg.ScriptExecFinished{Err: execErr}
					})
				}
				return m, cmds.RunScriptCapture(m.scriptService, item.ID)
			case entities.Package:
				m.state = ViewInstalling
				m.installStatus = "preparing"
				m.installMessage = "Preparando instalação..."
				m.installPercent = 0
				m.canAbort = false
				m.aborted = false
				return m, cmds.InstallPackage(m.installerService, item.ID)
			}
		} else {
			if m.confirmReturnOK {
				m.state = m.confirmReturn
			} else {
				m.state = ViewMainMenu
			}
			m.confirmReturnOK = false
		}
	}

	return m, nil
}
