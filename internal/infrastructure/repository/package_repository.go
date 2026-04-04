package repository

import (
	"sync"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// InMemoryPackageRepository is an in-memory implementation of PackageRepository
type InMemoryPackageRepository struct {
	packages map[string]*entities.Package
	mu       sync.RWMutex
}

// NewInMemoryPackageRepository creates a new in-memory package repository
func NewInMemoryPackageRepository() interfaces.PackageRepository {
	repo := &InMemoryPackageRepository{
		packages: make(map[string]*entities.Package),
	}
	repo.initializeDefaultPackages()
	return repo
}

// initializeDefaultPackages initializes the repository with default packages
func (r *InMemoryPackageRepository) initializeDefaultPackages() {
	defaultPackages := []entities.Package{
		// IDEs / Dev IDEs
		{
			ID:          "claude-code",
			Name:        "Claude Code CLI",
			Description: "CLI oficial da Anthropic para desenvolvimento com Claude",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://storage.googleapis.com/claude-code/install.sh",
			InstallCmd:  "bash install.sh",
			CheckCmd:    "which claude-code",
			ProjectURL:  "https://github.com/anthropics/claude-code",
		},
		{
			ID:          "cursor",
			Name:        "Cursor AI",
			Description: "Editor de código com IA integrada",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://download.cursor.sh/linux/appImage/x64",
			InstallCmd:  "chmod +x cursor.AppImage && sudo mv cursor.AppImage /usr/local/bin/cursor",
			CheckCmd:    "which cursor",
			ProjectURL:  "https://cursor.com",
		},
		{
			ID:          "antigravity",
			Name:        "Antigravity",
			Description: "IDE moderna com recursos avançados",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://antigravity.dev/download/linux",
			InstallCmd:  "sudo dpkg -i antigravity.deb || sudo apt-get install -f -y",
			CheckCmd:    "which antigravity",
			ProjectURL:  "https://antigravity.dev",
		},
		{
			ID:          "vscode",
			Name:        "VS Code",
			Description: "Editor de código da Microsoft, leve e extensível",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64",
			InstallCmd:  "sudo dpkg -i {{download_path}} || sudo apt-get install -f -y",
			CheckCmd:    "which code",
			ProjectURL:  "https://github.com/microsoft/vscode",
		},
		{
			ID:          "zed",
			Name:        "Zed",
			Description: "Editor de código moderno e rápido, focado em colaboração",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			DownloadURL: "https://zed.dev/api/releases/latest/linux/deb",
			InstallCmd:  "sudo dpkg -i {{download_path}} || sudo apt-get install -f -y",
			CheckCmd:    "which zed",
			ProjectURL:  "https://github.com/zed-industries/zed",
		},
		{
			ID:          "neovim",
			Name:        "Neovim",
			Description: "Fork moderno do Vim, extensível e otimizado para uso em terminal",
			Version:     "latest",
			Category:    types.PackageCategoryIDE,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y neovim",
			CheckCmd:    "which nvim",
			ProjectURL:  "https://github.com/neovim/neovim",
		},

		// Emuladores de Terminal
		{
			ID:          "wezterm",
			Name:        "WezTerm",
			Description: "Emulador de terminal altamente configurável, com renderização via GPU e multiplexação integrada",
			Version:     "latest",
			Category:    types.PackageCategoryTerminal,
			DownloadURL: "https://wezfurlong.org/wezterm/",
			InstallCmd:  "xdg-open https://wezfurlong.org/wezterm/ || echo \"Abra o site do WezTerm no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "which wezterm",
			ProjectURL:  "https://github.com/wez/wezterm",
		},
		{
			ID:          "kitty",
			Name:        "Kitty",
			Description: "Emulador de terminal rápido com renderização via GPU, suporte a imagens e multiplexação nativa",
			Version:     "latest",
			Category:    types.PackageCategoryTerminal,
			DownloadURL: "https://sw.kovidgoyal.net/kitty/",
			InstallCmd:  "xdg-open https://sw.kovidgoyal.net/kitty/ || echo \"Abra o site do Kitty no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "which kitty",
			ProjectURL:  "https://github.com/kovidgoyal/kitty",
		},
		{
			ID:          "alacritty",
			Name:        "Alacritty",
			Description: "Emulador de terminal minimalista e extremamente rápido, configurado via YAML",
			Version:     "latest",
			Category:    types.PackageCategoryTerminal,
			DownloadURL: "https://alacritty.org",
			InstallCmd:  "xdg-open https://alacritty.org || echo \"Abra o site do Alacritty no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "which alacritty",
			ProjectURL:  "https://github.com/alacritty/alacritty",
		},
		{
			ID:          "warp",
			Name:        "Warp",
			Description: "Terminal moderno com IA integrada, interface visual inovadora e blocos de comandos",
			Version:     "latest",
			Category:    types.PackageCategoryTerminal,
			DownloadURL: "https://warp.dev",
			InstallCmd:  "xdg-open https://warp.dev || echo \"Abra o site do Warp no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "which warp",
			ProjectURL:  "https://www.warp.dev",
		},
		{
			ID:          "wave-terminal",
			Name:        "Wave Terminal",
			Description: "Terminal moderno open-source inspirado no Warp, com foco em colaboração",
			Version:     "latest",
			Category:    types.PackageCategoryTerminal,
			DownloadURL: "https://waveterm.dev",
			InstallCmd:  "xdg-open https://waveterm.dev || echo \"Abra o site do Wave Terminal no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "which wave",
			ProjectURL:  "https://github.com/waveterm/waveterm",
		},

		// Shell Core (Zsh, Oh My Zsh, Powerlevel10k) - install via "Instalar componentes core"
		{
			ID:          "zsh",
			Name:        "Zsh",
			Description: "Z Shell - shell poderoso e configurável",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			InstallCmd:  "sudo apt-get install -y zsh",
			CheckCmd:    "which zsh",
			ProjectURL:  "https://sourceforge.net/projects/zsh/",
		},
		{
			ID:          "oh-my-zsh",
			Name:        "Oh My Zsh",
			Description: "Framework para gerenciar configuração Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			DownloadURL: "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh",
			InstallCmd:  "sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" \"\" --unattended",
			CheckCmd:    "test -d ~/.oh-my-zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh",
		},
		{
			ID:          "powerlevel10k",
			Name:        "Powerlevel10k",
			Description: "Tema Zsh rápido e customizável",
			Version:     "latest",
			Category:    types.PackageCategoryZshCore,
			InstallCmd:  "git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k",
			ProjectURL:  "https://github.com/romkatv/powerlevel10k",
		},

		// Shells alternativos
		{
			ID:          "fish-shell",
			Name:        "Fish Shell",
			Description: "Shell amigável e moderno, com autosuggestions e syntax highlighting nativos",
			Version:     "latest",
			Category:    types.PackageCategoryShell,
			InstallCmd:  "sudo apt-get install -y fish",
			CheckCmd:    "which fish",
			ProjectURL:  "https://github.com/fish-shell/fish-shell",
		},
		{
			ID:          "fisher",
			Name:        "Fisher",
			Description: "Gerenciador de plugins para Fish Shell",
			Version:     "latest",
			Category:    types.PackageCategoryShell,
			InstallCmd:  "fish -c 'curl -sL https://raw.githubusercontent.com/jorgebucaran/fisher/main/functions/fisher.fish | source && fisher install jorgebucaran/fisher'",
			CheckCmd:    "fish -c 'type -q fisher'",
			ProjectURL:  "https://github.com/jorgebucaran/fisher",
		},

		// Zsh Plugins - Built-in (5)
		{
			ID:          "zsh-plugin-git",
			Name:        "Git Plugin",
			Description: "Plugin built-in do Oh My Zsh para Git",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/git/git.plugin.zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/git",
		},
		{
			ID:          "zsh-plugin-docker",
			Name:        "Docker Plugin",
			Description: "Plugin built-in do Oh My Zsh para Docker",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/docker/docker.plugin.zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/docker",
		},
		{
			ID:          "zsh-plugin-rails",
			Name:        "Rails Plugin",
			Description: "Plugin built-in do Oh My Zsh para Ruby on Rails",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/rails/rails.plugin.zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/rails",
		},
		{
			ID:          "zsh-plugin-z",
			Name:        "Z Plugin",
			Description: "Plugin built-in para navegação rápida de diretórios",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/z/z.plugin.zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/z",
		},
		{
			ID:          "zsh-plugin-sudo",
			Name:        "Sudo Plugin",
			Description: "Plugin built-in para adicionar sudo facilmente",
			Version:     "built-in",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "# Built-in plugin",
			CheckCmd:    "test -f ~/.oh-my-zsh/plugins/sudo/sudo.plugin.zsh",
			ProjectURL:  "https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/sudo",
		},

		// Zsh Plugins - External (10)
		{
			ID:          "zsh-autosuggestions",
			Name:        "Zsh Autosuggestions",
			Description: "Sugestões automáticas baseadas no histórico",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			ProjectURL:  "https://github.com/zsh-users/zsh-autosuggestions",
		},
		{
			ID:          "zsh-syntax-highlighting",
			Name:        "Zsh Syntax Highlighting",
			Description: "Destaque de sintaxe para comandos Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting",
			ProjectURL:  "https://github.com/zsh-users/zsh-syntax-highlighting",
		},
		{
			ID:          "fzf-zsh",
			Name:        "FZF Zsh Integration",
			Description: "Integração do FZF com Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install --all",
			CheckCmd:    "test -d ~/.fzf",
			ProjectURL:  "https://github.com/junegunn/fzf",
		},
		{
			ID:          "you-should-use",
			Name:        "You Should Use",
			Description: "Lembra aliases existentes ao digitar comandos",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/MichaelAquilina/zsh-you-should-use.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/you-should-use",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/you-should-use",
			ProjectURL:  "https://github.com/MichaelAquilina/zsh-you-should-use",
		},
		{
			ID:          "zsh-completions",
			Name:        "Zsh Completions",
			Description: "Completions adicionais para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-completions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions",
			ProjectURL:  "https://github.com/zsh-users/zsh-completions",
		},
		{
			ID:          "zsh-history-substring-search",
			Name:        "Zsh History Substring Search",
			Description: "Busca no histórico por substring",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zsh-users/zsh-history-substring-search ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-history-substring-search",
			ProjectURL:  "https://github.com/zsh-users/zsh-history-substring-search",
		},
		{
			ID:          "fast-syntax-highlighting",
			Name:        "Fast Syntax Highlighting",
			Description: "Syntax highlighting mais rápido",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/zdharma-continuum/fast-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/fast-syntax-highlighting",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/fast-syntax-highlighting",
			ProjectURL:  "https://github.com/zdharma-continuum/fast-syntax-highlighting",
		},
		{
			ID:          "zsh-autocomplete",
			Name:        "Zsh Autocomplete",
			Description: "Autocomplete em tempo real para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone --depth 1 -- https://github.com/marlonrichert/zsh-autocomplete.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autocomplete",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autocomplete",
			ProjectURL:  "https://github.com/marlonrichert/zsh-autocomplete",
		},
		{
			ID:          "auto-notify",
			Name:        "Auto Notify",
			Description: "Notificações automáticas para comandos longos",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/MichaelAquilina/zsh-auto-notify.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/auto-notify",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/auto-notify",
			ProjectURL:  "https://github.com/MichaelAquilina/zsh-auto-notify",
		},
		{
			ID:          "zsh-vi-mode",
			Name:        "Zsh Vi Mode",
			Description: "Melhor modo Vi para Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "git clone https://github.com/jeffreytse/zsh-vi-mode ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-vi-mode",
			CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-vi-mode",
			ProjectURL:  "https://github.com/jeffreytse/zsh-vi-mode",
		},

		// Integração com IA
		{
			ID:          "shell-gpt",
			Name:        "ShellGPT",
			Description: "Assistente de IA universal para qualquer shell (explicações e sugestões inteligentes)",
			Version:     "latest",
			Category:    types.PackageCategoryAI,
			DownloadURL: "https://github.com/TheR1D/shell_gpt",
			InstallCmd:  "xdg-open https://github.com/TheR1D/shell_gpt || echo \"Abra o repositório ShellGPT no navegador para seguir as instruções de instalação (requer API key).\"",
			CheckCmd:    "command -v sgpt",
			ProjectURL:  "https://github.com/TheR1D/shell_gpt",
		},
		{
			ID:          "fish-ai",
			Name:        "Fish-AI",
			Description: "Integração de IA específica para Fish Shell, com sugestões inline",
			Version:     "latest",
			Category:    types.PackageCategoryAI,
			DownloadURL: "https://github.com/Realiserad/fish-ai",
			InstallCmd:  "xdg-open https://github.com/Realiserad/fish-ai || echo \"Abra o repositório fish-ai no navegador para seguir as instruções de instalação.\"",
			CheckCmd:    "test -d ~/.config/fish",
			ProjectURL:  "https://github.com/Realiserad/fish-ai",
		},

		// Development Tools (8)
		{
			ID:          "nvm",
			Name:        "NVM (Node Version Manager)",
			Description: "Gerenciador de versões Node.js",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash",
			CheckCmd:    "test -d ~/.nvm",
			ProjectURL:  "https://github.com/nvm-sh/nvm",
		},
		{
			ID:          "bun",
			Name:        "Bun",
			Description: "Runtime JavaScript/TypeScript rápido",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://bun.sh/install | bash",
			CheckCmd:    "test -d ~/.bun",
			ProjectURL:  "https://github.com/oven-sh/bun",
		},
		{
			ID:          "sdkman",
			Name:        "SDKMAN!",
			Description: "Gerenciador de SDKs para JVM",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -s https://get.sdkman.io | bash",
			CheckCmd:    "test -d ~/.sdkman",
			ProjectURL:  "https://github.com/sdkman/sdkman-cli",
		},
		{
			ID:          "pnpm",
			Name:        "pnpm",
			Description: "Gerenciador de pacotes Node.js eficiente",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://get.pnpm.io/install.sh | sh -",
			CheckCmd:    "which pnpm",
			ProjectURL:  "https://github.com/pnpm/pnpm",
		},
		{
			ID:          "deno",
			Name:        "Deno",
			Description: "Runtime seguro para JavaScript e TypeScript",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://deno.land/install.sh | sh",
			CheckCmd:    "which deno",
			ProjectURL:  "https://github.com/denoland/deno",
		},
		{
			ID:          "angular-cli",
			Name:        "Angular CLI",
			Description: "Interface de linha de comando para Angular",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "npm install -g @angular/cli",
			CheckCmd:    "which ng",
			ProjectURL:  "https://github.com/angular/angular-cli",
		},
		{
			ID:          "openvpn3",
			Name:        "OpenVPN 3",
			Description: "Cliente VPN moderno",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "sudo apt-get install -y openvpn3",
			CheckCmd:    "which openvpn3",
			ProjectURL:  "https://community.openvpn.net/openvpn/wiki/OpenVPN3Linux",
		},
		{
			ID:          "gh",
			Name:        "GitHub CLI",
			Description: "CLI oficial do GitHub; necessário para criar repositórios automaticamente no fluxo Configurar Zsh",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main\" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null && sudo apt update && sudo apt install gh -y",
			CheckCmd:    "which gh",
			ProjectURL:  "https://github.com/cli/cli",
		},
		{
			ID:          "homebrew",
			Name:        "Homebrew",
			Description: "Gerenciador de pacotes para Linux",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"",
			CheckCmd:    "which brew",
			ProjectURL:  "https://github.com/Homebrew/install",
		},
		{
			ID:          "openjdk",
			Name:        "Java OpenJDK",
			Description: "Implementação open-source da plataforma Java (JDK)",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y openjdk-21-jdk",
			CheckCmd:    "which javac",
			ProjectURL:  "https://openjdk.org/",
		},
		{
			ID:          "starship",
			Name:        "Starship",
			Description: "Prompt minimalista, rápido e personalizável para qualquer shell",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -sS https://starship.rs/install.sh | sh -s -- -y",
			CheckCmd:    "which starship",
			ProjectURL:  "https://github.com/starship/starship",
		},
		{
			ID:          "flathub",
			Name:        "Flathub (Flatpak Remote)",
			Description: "Repositório principal de aplicações Flatpak",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo",
			CheckCmd:    "flatpak remote-list | grep -q flathub",
			ProjectURL:  "https://github.com/flathub/flathub",
		},
		{
			ID:          "google-chrome",
			Name:        "Google Chrome",
			Description: "Navegador web Google Chrome (versão estável)",
			Version:     "latest",
			Category:    types.PackageCategoryApp,
			DownloadURL: "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb",
			InstallCmd:  "sudo dpkg -i {{download_path}} || sudo apt-get install -f -y",
			CheckCmd:    "which google-chrome || which google-chrome-stable",
			ProjectURL:  "https://www.google.com/chrome/",
		},
		{
			ID:          "insomnia",
			Name:        "Insomnia",
			Description: "Cliente HTTP e GraphQL para testar APIs",
			Version:     "latest",
			Category:    types.PackageCategoryApp,
			InstallCmd:  "flatpak install -y flathub rest.insomnia.Insomnia",
			CheckCmd:    "flatpak list | grep -q rest.insomnia.Insomnia",
			ProjectURL:  "https://github.com/Kong/insomnia",
		},
		{
			ID:          "remmina",
			Name:        "Remmina",
			Description: "Cliente de desktop remoto com suporte a RDP, VNC, SPICE, X2Go, SSH e mais",
			Version:     "latest",
			Category:    types.PackageCategoryApp,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y remmina",
			CheckCmd:    "which remmina",
			ProjectURL:  "https://gitlab.com/Remmina/Remmina",
		},
		{
			ID:          "distrobox",
			Name:        "Distrobox",
			Description: "Contêineres integrados ao sistema para múltiplas distribuições Linux",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl -s https://raw.githubusercontent.com/89luca89/distrobox/main/install | sudo sh",
			CheckCmd:    "which distrobox",
			ProjectURL:  "https://github.com/89luca89/distrobox",
		},
		{
			ID:          "mise",
			Name:        "Mise",
			Description: "Gerenciador de versões para múltiplas linguagens e ferramentas",
			Version:     "latest",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "curl https://mise.run | sh",
			CheckCmd:    "which mise",
			ProjectURL:  "https://github.com/jdx/mise",
		},
		{
			ID:          "dotnet-sdk",
			Name:        ".NET SDK",
			Description: "SDK do .NET para desenvolvimento de aplicações",
			Version:     "8.0",
			Category:    types.PackageCategoryTool,
			InstallCmd:  "wget https://dot.net/v1/dotnet-install.sh -O /tmp/dotnet-install.sh && chmod +x /tmp/dotnet-install.sh && /tmp/dotnet-install.sh --version latest && echo 'export DOTNET_ROOT=$HOME/.dotnet' >> ~/.bashrc && echo 'export PATH=$PATH:$HOME/.dotnet:$HOME/.dotnet/tools' >> ~/.bashrc",
			CheckCmd:    "which dotnet",
			ProjectURL:  "https://github.com/dotnet/sdk",
		},

		// Administração de sistemas
		{
			ID:          "cockpit",
			Name:        "Cockpit (servidor)",
			Description: "Painel web para administrar esta máquina (serviços, logs, armazenamento, rede)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y cockpit && sudo systemctl enable --now cockpit.socket",
			CheckCmd:    "systemctl is-active cockpit.socket >/dev/null 2>&1 || command -v cockpit-bridge >/dev/null 2>&1",
			ProjectURL:  "https://github.com/cockpit-project/cockpit",
			Notes: "Fluxo pensado para Ubuntu/Debian (apt). Expõe interface web (porta 9090 por padrão).\n\nApós instalar: acesse https://IP-DESTA-MAQUINA:9090 no navegador.\n\nSe usar firewalld (Fedora/RHEL): sudo firewall-cmd --add-service=cockpit --permanent && sudo firewall-cmd --reload\n\nEm Debian estável mais antigo o pacote pode estar só em backports — ajuste sources.list se o apt não achar cockpit.\n\nDocumentação: https://cockpit-project.org/",
		},
		{
			ID:          "cockpit-client",
			Name:        "Cockpit Client",
			Description: "App para conectar a servidores Cockpit remotos (Flatpak)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "flatpak install -y flathub org.cockpit_project.CockpitClient",
			CheckCmd:    "flatpak list --app 2>/dev/null | grep -q org.cockpit_project.CockpitClient",
			ProjectURL:  "https://github.com/cockpit-project/cockpit",
			Notes: "Requer Flatpak e o remote flathub (use o instalador Flathub em Ferramentas se precisar).\n\nÚtil quando o servidor já roda Cockpit e você só quer o cliente nesta máquina.\n\nDocumentação: https://cockpit-project.org/",
		},
		{
			ID:          "webmin",
			Name:        "Webmin",
			Description: "Painel web clássico de administração Unix/Linux",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "curl -fsSL -o /tmp/webmin-setup-repo.sh https://raw.githubusercontent.com/webmin/webmin/master/webmin-setup-repo.sh && sudo sh /tmp/webmin-setup-repo.sh -f && rm -f /tmp/webmin-setup-repo.sh && sudo apt-get update && sudo apt-get install -y webmin",
			CheckCmd:    "dpkg -l webmin 2>/dev/null | grep -q '^ii'",
			ProjectURL:  "https://github.com/webmin/webmin",
			Notes: "Adiciona repositório oficial via script da Webmin e instala o pacote (apt).\n\nServiço web com privilégios elevados: use senha forte, firewall e rede confiável.\n\nApós instalar, o script costuma indicar URL e porta (muitas vezes https://localhost:10000).\n\nDocumentação: https://webmin.com/",
		},
		{
			ID:          "topgrade",
			Name:        "Topgrade",
			Description: "Atualiza pacotes do sistema, flatpak, rustup, npm global e dezenas de outras fontes num único comando",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y pipx && pipx install topgrade",
			CheckCmd:    "test -x \"$HOME/.local/bin/topgrade\"",
			ProjectURL:  "https://github.com/topgrade-rs/topgrade",
			Notes: "Instalação via pipx no seu usuário (sem misturar com apt do sistema).\n\nSe `topgrade` não aparecer no PATH, abra um novo terminal ou execute: export PATH=\"$HOME/.local/bin:$PATH\"\n\nNa primeira execução o Topgrade pode pedir confirmações — leia cada etapa.\n\nDocumentação: https://github.com/topgrade-rs/topgrade",
		},
		{
			ID:          "termius",
			Name:        "Termius",
			Description: "Cliente SSH/SFTP com sincronização (Flatpak)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "flatpak install -y flathub com.termius.Termius",
			CheckCmd:    "flatpak list --app 2>/dev/null | grep -q com.termius.Termius",
			ProjectURL:  "https://termius.com",
			Notes: "Requer Flatpak + flathub.\n\nConta Termius e dados na nuvem são opcionais; revise a política de privacidade se for usar sync.\n\nSite: https://termius.com/",
		},
		{
			ID:          "cpu-x",
			Name:        "CPU-X",
			Description: "Informações detalhadas de CPU, GPU, placa-mãe e sensores (Flatpak)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "flatpak install -y flathub io.github.thetumultuousunicornofdarkness.cpu-x",
			CheckCmd:    "flatpak list --app 2>/dev/null | grep -q io.github.thetumultuousunicornofdarkness.cpu-x",
			ProjectURL:  "https://github.com/TheTumultuousUnicornOfDarkness/CPU-X",
			Notes: "Requer Flatpak + flathub.\n\nDocumentação: https://thetumultuousunicornofdarkness.github.io/CPU-X/",
		},
		{
			ID:          "sssd-active-directory",
			Name:        "Ferramentas Active Directory (realm/sssd)",
			Description: "Pacotes para integrar Ubuntu/Debian a um domínio AD (SSSD, realmd, Kerberos)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "sudo apt-get update && sudo apt-get install -y sssd realmd adcli samba-common-bin adsys krb5-user libpam-krb5 libpam-ccreds auth-client-config oddjob oddjob-mkhomedir",
			CheckCmd:    "command -v realm >/dev/null 2>&1 && dpkg -l sssd 2>/dev/null | grep -q '^ii'",
			ProjectURL:  "https://github.com/SSSD/sssd",
			Notes: "Este passo só INSTALA pacotes; não junta o computador ao domínio sozinho.\n\nDepois você precisa configurar DNS, horário (NTP) e executar algo como: sudo realm join DOMINIO.EXEMPLO -U administrador\n\nSe o apt reclamar do pacote adsys (comum em algumas versões do Debian), instale os demais pacotes manualmente sem o adsys.\n\nRevise com seu time de infraestrutura antes de alterar autenticação em produção.\n\nGuia Ubuntu: https://ubuntu.com/server/docs/samba-ad-integration",
		},
		{
			ID:          "sloth-bash",
			Name:        "Sloth-Bash",
			Description: "Ferramentas shell do projeto Sloth-Bash (script remoto oficial)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "curl -fsSL https://raw.githubusercontent.com/psygreg/sloth-bash/main/install.sh | bash",
			CheckCmd:    "command -v sloth >/dev/null 2>&1 || test -d \"$HOME/.sloth-bash\"",
			ProjectURL:  "https://github.com/psygreg/sloth-bash",
			Notes: "Executa o install.sh oficial do repositório (curl | bash). Só confirme se você confia na origem.\n\nRepositório: https://github.com/psygreg/sloth-bash",
		},
		{
			ID:          "zash-terminal",
			Name:        "Zash Terminal",
			Description: "Terminal com recursos extras (instalador upstream)",
			Version:     "latest",
			Category:    types.PackageCategorySysAdmin,
			InstallCmd:  "curl -fsSL https://raw.githubusercontent.com/leoberbert/zashterminal/refs/heads/main/install.sh -o /tmp/zashterminal-install.sh && chmod +x /tmp/zashterminal-install.sh && bash /tmp/zashterminal-install.sh && rm -f /tmp/zashterminal-install.sh",
			CheckCmd:    "command -v zashterminal >/dev/null 2>&1",
			ProjectURL:  "https://github.com/leoberbert/zashterminal",
			Notes: "O script oficial pode pedir sudo e fazer alterações no sistema.\n\nLeia o que o instalador imprimir antes de aceitar.\n\nRepositório: https://github.com/leoberbert/zashterminal",
		},

		// Games
		{
			ID:          "prism-launcher",
			Name:        "Prism Launcher",
			Description: "Launcher open-source para Minecraft com múltiplas instâncias",
			Version:     "latest",
			Category:    types.PackageCategoryGames,
			InstallCmd:  "flatpak install -y flathub org.prismlauncher.PrismLauncher",
			CheckCmd:    "flatpak list | grep -q org.prismlauncher.PrismLauncher",
			ProjectURL:  "https://github.com/PrismLauncher/PrismLauncher",
		},
		{
			ID:          "lutris",
			Name:        "Lutris",
			Description: "Plataforma de jogos para Linux, integrando Wine, emuladores e lojas",
			Version:     "latest",
			Category:    types.PackageCategoryGames,
			InstallCmd:  "flatpak install -y flathub net.lutris.Lutris",
			CheckCmd:    "flatpak list | grep -q net.lutris.Lutris",
			ProjectURL:  "https://github.com/lutris/lutris",
		},
		{
			ID:          "gear-lever",
			Name:        "Gear Lever",
			Description: "Gerenciador gráfico para integrar e atualizar AppImages no sistema",
			Version:     "latest",
			Category:    types.PackageCategoryApp,
			InstallCmd:  "flatpak install -y flathub it.mijorus.gearlever",
			CheckCmd:    "flatpak list | grep -q it.mijorus.gearlever",
			ProjectURL:  "https://github.com/mijorus/gearlever",
		},
	}

	for _, pkg := range defaultPackages {
		pkgCopy := pkg
		r.packages[pkg.ID] = &pkgCopy
	}
}

// FindAll returns all packages
func (r *InMemoryPackageRepository) FindAll() ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0, len(r.packages))
	for _, pkg := range r.packages {
		packages = append(packages, *pkg)
	}

	return packages, nil
}

// FindByID finds a package by ID
func (r *InMemoryPackageRepository) FindByID(id string) (*entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pkg, exists := r.packages[id]
	if !exists {
		return nil, types.ErrNotFound
	}

	pkgCopy := *pkg
	return &pkgCopy, nil
}

// FindByCategory finds packages by category
func (r *InMemoryPackageRepository) FindByCategory(category types.PackageCategory) ([]entities.Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packages := make([]entities.Package, 0)
	for _, pkg := range r.packages {
		if pkg.Category == category {
			packages = append(packages, *pkg)
		}
	}

	return packages, nil
}

// Save saves a package
func (r *InMemoryPackageRepository) Save(pkg *entities.Package) error {
	if err := pkg.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	pkgCopy := *pkg
	r.packages[pkg.ID] = &pkgCopy

	return nil
}

// Delete deletes a package
func (r *InMemoryPackageRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.packages[id]; !exists {
		return types.ErrNotFound
	}

	delete(r.packages, id)
	return nil
}

// Exists checks if a package exists
func (r *InMemoryPackageRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.packages[id]
	return exists
}
