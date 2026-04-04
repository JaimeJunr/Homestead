// Package items implements bubble list.Item for the main TUI lists.
package items

import (
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

type MenuItem struct {
	Label  string
	Desc   string
	Action string
}

func (i MenuItem) Title() string       { return i.Label }
func (i MenuItem) Description() string { return i.Desc }
func (i MenuItem) FilterValue() string { return i.Label }

type ScriptItem struct {
	Script entities.Script
}

func (i ScriptItem) Title() string       { return i.Script.Name }
func (i ScriptItem) Description() string { return i.Script.Description }
func (i ScriptItem) FilterValue() string { return i.Script.Name }

type PackageItem struct {
	Pkg entities.Package
}

func (i PackageItem) Title() string       { return i.Pkg.Name }
func (i PackageItem) Description() string { return i.Pkg.Description }
func (i PackageItem) FilterValue() string { return i.Pkg.Name }

type InstallerCategoryItem struct {
	Heading    string
	Desc       string
	Categories []types.PackageCategory
	Utilities  bool // open utilities scripts under Instaladores
}

func (i InstallerCategoryItem) Title() string       { return i.Heading }
func (i InstallerCategoryItem) Description() string { return i.Desc }
func (i InstallerCategoryItem) FilterValue() string { return i.Heading }
