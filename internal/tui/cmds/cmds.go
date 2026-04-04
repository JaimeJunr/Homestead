// Package cmds builds tea.Cmd values for async work (catalog, install, scripts, URLs).
package cmds

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
	"github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/sysurl"
)

// CheckZshCoreInstalled detects oh-my-zsh for the main menu.
func CheckZshCoreInstalled(installerService *services.InstallerService) tea.Cmd {
	return func() tea.Msg {
		installed, _ := installerService.IsPackageInstalled("oh-my-zsh")
		return msg.ZshCoreInstalled{Installed: installed}
	}
}

// FetchCatalog downloads and merges the remote installer catalog when url is non-empty.
func FetchCatalog(url string, svc *services.InstallerService) tea.Cmd {
	if strings.TrimSpace(url) == "" {
		return nil
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		body, err := catalog.Fetch(ctx, url)
		if err != nil {
			return msg.CatalogFetched{Err: err}
		}
		pkgs, _, err := catalog.ParseManifest(body)
		if err != nil {
			return msg.CatalogFetched{Err: err}
		}
		if err := svc.MergePackages(pkgs); err != nil {
			return msg.CatalogFetched{Err: err}
		}
		_ = catalog.WriteCache(body)
		return msg.CatalogFetched{Ok: true}
	}
}

// RunScriptCapture runs ExecuteScriptCapture in a Cmd.
func RunScriptCapture(service *services.ScriptService, scriptID string) tea.Cmd {
	return func() tea.Msg {
		out, err := service.ExecuteScriptCapture(scriptID)
		return msg.ScriptCaptured{Output: out, Err: err}
	}
}

// InstallPackage streams install progress as Progress messages.
func InstallPackage(service *services.InstallerService, packageID string) tea.Cmd {
	return func() tea.Msg {
		progressChan := make(chan interfaces.InstallProgress, 10)

		go func() {
			err := service.InstallPackage(packageID, func(progress interfaces.InstallProgress) {
				progressChan <- progress
			})
			if err != nil {
				progressChan <- interfaces.InstallProgress{
					Status:      "failed",
					Message:     err.Error(),
					IsCompleted: true,
					Error:       err,
				}
			}
			close(progressChan)
		}()

		for progress := range progressChan {
			return msg.Progress(progress)
		}

		return msg.InstallComplete{Err: nil}
	}
}

// ApplyZshConfig runs ConfigService.ApplyConfig and sends ZshApplyResult.
func ApplyZshConfig(configService *services.ConfigService, selections interfaces.ConfigSelections) tea.Cmd {
	return func() tea.Msg {
		err := configService.ApplyConfig(selections)
		return msg.ZshApplyResult{Err: err}
	}
}

// OpenURL opens url and reports with URLActionDone (Verb "open").
func OpenURL(url string) tea.Cmd {
	return func() tea.Msg {
		err := sysurl.Open(url)
		return msg.URLActionDone{Verb: "open", Err: err}
	}
}

// CopyURL copies url to clipboard and reports with URLActionDone (Verb "copy").
func CopyURL(url string) tea.Cmd {
	return func() tea.Msg {
		err := sysurl.CopyToClipboard(url)
		return msg.URLActionDone{Verb: "copy", Err: err}
	}
}
