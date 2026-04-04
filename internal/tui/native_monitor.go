package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
)

const nativeMonitorRefreshInterval = 3 * time.Second

type nativeMonitorReloadMsg struct {
	kind    string
	battery *monitoring.BatterySnapshot
	memory  *monitoring.MemorySnapshot
	err     error
}

type nativeMonitorTickMsg struct{}

func nativeMonitorScheduleTick() tea.Cmd {
	return tea.Tick(nativeMonitorRefreshInterval, func(time.Time) tea.Msg {
		return nativeMonitorTickMsg{}
	})
}

func (m Model) nativeMonitorLoadCmd() tea.Cmd {
	kind := m.nativeMonitorKind
	return func() tea.Msg {
		msg := nativeMonitorReloadMsg{kind: kind}
		switch kind {
		case entities.NativeMonitorBattery:
			msg.battery, msg.err = monitoring.ReadBattery()
		case entities.NativeMonitorMemory:
			msg.memory, msg.err = monitoring.ReadMemory()
		default:
			msg.err = fmt.Errorf("monitor desconhecido: %q", kind)
		}
		return msg
	}
}

func (m Model) renderNativeMonitorView() string {
	boxW := scriptOutputCardWidth(m.width)
	head := titleStyle.Render("Homestead") + "\n" +
		helpStyle.Render("Gerenciador de Sistema") + "\n" +
		scriptOutputDivider(boxW) + "\n"

	var body string
	switch m.nativeMonitorKind {
	case entities.NativeMonitorBattery:
		body = renderNativeBatteryPanel(m)
	case entities.NativeMonitorMemory:
		body = renderNativeMemoryPanel(m)
	default:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Monitor inválido.")
	}

	footer := scriptScreenFooterBarStyle.Width(max(12, boxW-8)).Render(
		"r: atualizar agora · Enter / Esc / q: voltar · atualiza a cada 3s",
	)
	content := head + body + "\n" + footer
	box := scriptScreenOuterStyle.Width(boxW)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

func renderNativeBatteryPanel(m Model) string {
	title := scriptScreenAccentStyle.Render("🔋 Monitor de bateria")

	if m.nativeBatteryErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeBatteryErr.Error())
	}
	b := m.nativeBattery
	if b == nil {
		return title + "\n\n" + helpStyle.Render("Carregando…")
	}

	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")

	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(22).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}

	status := b.Status
	if status == "Charging" {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(status)
	} else if status == "Discharging" {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(status)
	} else {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(status)
	}
	kv("Status", status)

	if b.Capacity >= 0 {
		kv("Capacidade", fmt.Sprintf("%d %%", b.Capacity))
	}
	if b.CapacityLevel != "" && b.CapacityLevel != "Unknown" {
		kv("Nível", b.CapacityLevel)
	}

	if b.ACOnline != nil {
		var acLine string
		if *b.ACOnline {
			acLine = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("conectado")
		} else {
			acLine = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("desconectado")
		}
		if b.ACName != "" {
			acLine += helpStyle.Render("  (" + b.ACName + ")")
		}
		kv("Carregador", acLine)
	}

	sb.WriteString("\n")
	sb.WriteString(scriptScreenAccentStyle.Render("Detalhes") + "\n")

	if b.EnergyNowUWh > 0 {
		kv("Energia agora", fmt.Sprintf("%.2f Wh", float64(b.EnergyNowUWh)/1e6))
	}
	if b.EnergyFullUWh > 0 {
		kv("Energia cheia", fmt.Sprintf("%.2f Wh", float64(b.EnergyFullUWh)/1e6))
	}
	if b.EnergyDesignUWh > 0 {
		kv("Design (cheia)", fmt.Sprintf("%.2f Wh", float64(b.EnergyDesignUWh)/1e6))
	}
	if h, ok := b.HealthPercent(); ok {
		kv("Saúde (est.)", fmt.Sprintf("%.1f %%", h))
	}
	if w, ok := b.PowerWatts(); ok {
		prefix := ""
		if b.Status == "Charging" {
			prefix = "+"
		} else if b.Status == "Discharging" {
			prefix = "−"
		}
		kv("Potência", fmt.Sprintf("%s%.2f W", prefix, w))
	}
	if b.VoltageNowUV > 0 {
		kv("Tensão", fmt.Sprintf("%.2f V", float64(b.VoltageNowUV)/1e6))
	}
	if b.CycleCount > 0 {
		kv("Ciclos", fmt.Sprintf("%d", b.CycleCount))
	}
	if b.Technology != "" {
		kv("Tecnologia", b.Technology)
	}
	if b.Manufacturer != "" {
		kv("Fabricante", b.Manufacturer)
	}
	if b.ModelName != "" {
		kv("Modelo", b.ModelName)
	}

	sb.WriteString("\n")
	sb.WriteString(scriptScreenAccentStyle.Render("Resumo") + "\n")
	sb.WriteString(batteryStatusLine(b))
	if est := batteryChargeETA(b); est != "" {
		sb.WriteString(helpStyle.Render(est) + "\n")
	}

	return sb.String()
}

func batteryStatusLine(b *monitoring.BatterySnapshot) string {
	if b == nil {
		return ""
	}
	acOn := b.ACOnline != nil && *b.ACOnline
	switch {
	case b.Status == "Charging" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓ Carregando com o adaptador conectado.")
	case b.Status == "Full" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✓ Bateria cheia.")
	case b.Status == "Discharging" && !acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("⚠ Descarregando (sem AC).")
	case b.Status == "Discharging" && acOn:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("⚠ AC conectado, mas status é descarga — verifique cabo/adaptador.")
	default:
		return helpStyle.Render("Status: " + b.Status)
	}
}

func batteryChargeETA(b *monitoring.BatterySnapshot) string {
	if b == nil || b.Status != "Charging" || b.PowerNowUW <= 0 {
		return ""
	}
	if b.EnergyFullUWh <= 0 || b.EnergyNowUWh < 0 {
		return ""
	}
	rem := b.EnergyFullUWh - b.EnergyNowUWh
	if rem <= 0 {
		return ""
	}
	h := float64(rem) / float64(b.PowerNowUW)
	if h <= 0 || h > 48 {
		return ""
	}
	hh := int(h)
	mm := int((h - float64(hh)) * 60)
	return fmt.Sprintf("Tempo estimado até cheio: ~%dh %dm", hh, mm)
}

func renderNativeMemoryPanel(m Model) string {
	title := scriptScreenAccentStyle.Render("🧠 Uso de memória")

	if m.nativeMemoryErr != nil {
		return title + "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.nativeMemoryErr.Error())
	}
	s := m.nativeMemory
	if s == nil {
		return title + "\n\n" + helpStyle.Render("Carregando…")
	}

	mb := func(kb uint64) string {
		return fmt.Sprintf("%.0f MiB", float64(kb)/1024)
	}

	var sb strings.Builder
	sb.WriteString(title)
	sb.WriteString("\n\n")

	kv := func(k, v string) {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(14).Render(k))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v))
		sb.WriteString("\n")
	}

	sb.WriteString(scriptScreenAccentStyle.Render("RAM") + "\n")
	kv("Total", mb(s.MemTotalKB))
	kv("Usado*", mb(s.UsedApproxKB()))
	kv("Livre", mb(s.MemFreeKB))
	if s.MemAvailableKB > 0 {
		kv("Disponível", mb(s.MemAvailableKB))
	}
	kv("Compart.", mb(s.ShmemKB))

	sb.WriteString("\n")
	sb.WriteString(scriptScreenAccentStyle.Render("Swap") + "\n")
	kv("Total", mb(s.SwapTotalKB))
	kv("Livre", mb(s.SwapFreeKB))
	if s.SwapTotalKB > s.SwapFreeKB {
		kv("Usado", mb(s.SwapTotalKB-s.SwapFreeKB))
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("* “Usado” é uma estimativa."))

	return sb.String()
}
