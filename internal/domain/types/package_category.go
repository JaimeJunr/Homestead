package types

// PackageCategory represents the category of a package
type PackageCategory string

const (
	PackageCategoryIDE  PackageCategory = "ide"
	PackageCategoryTool PackageCategory = "tool"
	PackageCategoryApp  PackageCategory = "app"
)

// IsValid checks if the category is valid
func (c PackageCategory) IsValid() bool {
	switch c {
	case PackageCategoryIDE, PackageCategoryTool, PackageCategoryApp:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (c PackageCategory) String() string {
	return string(c)
}
