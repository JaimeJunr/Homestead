// Package tui is the Bubble Tea front-end for Homestead.
//
// Layout:
//
//   - Root (this package): Model, Update/View, navigation, list loading, Zsh wizards, native monitor.
//   - cmds: tea.Cmd factories (catalog fetch, install, script capture, URLs).
//   - items: bubble list.Item implementations (menu, scripts, packages, installer groups).
//   - msg: Bubble Tea message types for async and cross-screen signals.
//   - sysurl: open URL / clipboard without importing the root tui package (avoids cycles).
//   - theme: shared Lipgloss styles and installer section titles.
package tui
