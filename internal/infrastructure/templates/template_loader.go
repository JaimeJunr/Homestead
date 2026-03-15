package templates

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateLoader handles loading and rendering Go templates
type TemplateLoader struct {
	templateDir string
	embeddedFS  *embed.FS
	fsPrefix    string
	cache       map[string]*template.Template
}

// NewTemplateLoader creates a new template loader backed by the filesystem
func NewTemplateLoader(templateDir string) *TemplateLoader {
	return &TemplateLoader{
		templateDir: templateDir,
		cache:       make(map[string]*template.Template),
	}
}

// NewTemplateLoaderFromFS creates a loader backed by an embed.FS, falling back to disk if templateDir is set
func NewTemplateLoaderFromFS(fs embed.FS, prefix string) *TemplateLoader {
	return &TemplateLoader{
		embeddedFS: &fs,
		fsPrefix:   prefix,
		cache:      make(map[string]*template.Template),
	}
}

// LoadTemplate loads a template, checking embedded FS first then falling back to disk
func (tl *TemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	// Check cache first
	if tmpl, ok := tl.cache[name]; ok {
		return tmpl, nil
	}

	var content []byte

	// Try embedded FS first if available
	if tl.embeddedFS != nil {
		fsPath := name
		if tl.fsPrefix != "" {
			fsPath = tl.fsPrefix + "/" + name
		}
		var readErr error
		content, readErr = tl.embeddedFS.ReadFile(fsPath)
		if readErr != nil {
			content = nil
		}
	}

	// Fall back to filesystem if not loaded from embedded FS
	if content == nil && tl.templateDir != "" {
		templatePath := filepath.Join(tl.templateDir, name)
		var readErr error
		content, readErr = os.ReadFile(templatePath)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				return nil, fmt.Errorf("template '%s' not found", name)
			}
			return nil, fmt.Errorf("failed to read template '%s': %w", name, readErr)
		}
	}

	if content == nil {
		return nil, fmt.Errorf("template '%s' not found", name)
	}

	// Parse template
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template '%s': %w", name, err)
	}

	// Cache template
	tl.cache[name] = tmpl

	return tmpl, nil
}

// RenderTemplate loads and renders a template with the provided data
func (tl *TemplateLoader) RenderTemplate(name string, data interface{}) (string, error) {
	// Load template
	tmpl, err := tl.LoadTemplate(name)
	if err != nil {
		return "", err
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", name, err)
	}

	return buf.String(), nil
}

// ListTemplates returns a list of all available template files from the filesystem
func (tl *TemplateLoader) ListTemplates() ([]string, error) {
	if tl.templateDir == "" {
		return []string{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(tl.templateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read template directory: %w", err)
	}

	// Collect template files
	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only include .tmpl files
		if strings.HasSuffix(name, ".tmpl") {
			templates = append(templates, name)
		}
	}

	return templates, nil
}

// ClearCache clears the template cache
func (tl *TemplateLoader) ClearCache() {
	tl.cache = make(map[string]*template.Template)
}

// HasTemplate checks if a template exists (embedded FS first, then disk)
func (tl *TemplateLoader) HasTemplate(name string) bool {
	if tl.embeddedFS != nil {
		fsPath := name
		if tl.fsPrefix != "" {
			fsPath = tl.fsPrefix + "/" + name
		}
		f, err := tl.embeddedFS.Open(fsPath)
		if err == nil {
			f.Close()
			return true
		}
	}

	if tl.templateDir != "" {
		templatePath := filepath.Join(tl.templateDir, name)
		_, err := os.Stat(templatePath)
		return err == nil
	}

	return false
}
