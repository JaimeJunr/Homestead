package templates

import "embed"

//go:embed files/*.tmpl
var EmbeddedTemplates embed.FS
