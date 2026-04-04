package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
)

// OpenURL opens url with the OS default handler (xdg-open / gio on Linux).
func OpenURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("URL vazia")
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		if path, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command(path, url)
		} else if path, err := exec.LookPath("gio"); err == nil {
			cmd = exec.Command(path, "open", url)
		} else {
			return fmt.Errorf("instale xdg-utils (xdg-open) ou glib2 (gio) para abrir links no navegador")
		}
	}
	return cmd.Start()
}

// CopyURLToClipboard copies text to the clipboard (Wayland / X11).
func CopyURLToClipboard(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("nada para copiar")
	}
	if path, err := exec.LookPath("wl-copy"); err == nil {
		c := exec.Command(path)
		c.Stdin = strings.NewReader(text)
		return c.Run()
	}
	if path, err := exec.LookPath("xclip"); err == nil {
		c := exec.Command(path, "-selection", "clipboard")
		c.Stdin = strings.NewReader(text)
		return c.Run()
	}
	if path, err := exec.LookPath("xsel"); err == nil {
		c := exec.Command(path, "--clipboard", "--input")
		c.Stdin = strings.NewReader(text)
		return c.Run()
	}
	return fmt.Errorf("instale wl-copy (Wayland) ou xclip/xsel (X11) para copiar no terminal")
}

// PackageKeyboardURL is ProjectURL, or DownloadURL if ProjectURL is empty.
func PackageKeyboardURL(pkg entities.Package) string {
	if u := strings.TrimSpace(pkg.ProjectURL); u != "" {
		return u
	}
	if u := strings.TrimSpace(pkg.DownloadURL); u != "" {
		return u
	}
	return ""
}
