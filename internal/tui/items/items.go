// Package items provides bubble list.Item implementations for the main TUI.
package items

import (
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// MenuItem is a main menu row.
type MenuItem struct {
	Label  string
	Desc   string
	Action string
}

func (i MenuItem) Title() string       { return i.Label }
func (i MenuItem) Description() string { return i.Desc }
func (i MenuItem) FilterValue() string { return i.Label }

// ScriptItem wraps a script for the list.
type ScriptItem struct {
	Script entities.Script
}

func (i ScriptItem) Title() string       { return i.Script.Name }
func (i ScriptItem) Description() string { return i.Script.Description }
func (i ScriptItem) FilterValue() string { return i.Script.Name }

// PackageItem wraps a package for the list.
type PackageItem struct {
	Pkg entities.Package
}

func (i PackageItem) Title() string       { return i.Pkg.Name }
func (i PackageItem) Description() string { return i.Pkg.Description }
func (i PackageItem) FilterValue() string { return i.Pkg.Name }

// InstallerCategoryItem groups installer catalog categories (or utilities scripts).
type InstallerCategoryItem struct {
	Heading    string
	Desc       string
	Categories []types.PackageCategory
	Utilities  bool // true: open CategoryUtilities script list under Instaladores
}

func (i InstallerCategoryItem) Title() string       { return i.Heading }
func (i InstallerCategoryItem) Description() string { return i.Desc }
func (i InstallerCategoryItem) FilterValue() string { return i.Heading }
