// Package templateprovider defines interfaces and implementations for template
// engines used to render SKA blueprints (e.g., Sprig and Jinja-like templates).
package templateprovider

import (
	"io"
)

// TemplateService abstracts a template engine implementation.
// Implementations must load templates from files or strings and execute them
// with a set of variables. Optional helpers exist to detect missing keys and
// optional-field control flow.
type TemplateService interface {
	// FromFile loads a template from a file path.
	FromFile(path string) error
	// FromString loads a template from the given string content.
	FromString(templateContent string) error
	// Execute renders the template to the given writer using the provided variables.
	Execute(fp io.Writer, withVariables map[string]interface{}) error
	// WithErrorOnMissingKey toggles erroring on missing keys during execution, if supported.
	WithErrorOnMissingKey(key bool)
	// IsMissingKeyError reports whether the provided error indicates a missing key.
	IsMissingKeyError(err error) bool
	// IsOptionalError reports whether the provided error indicates an optional control flow skip.
	IsOptionalError(err error) bool
}

// TemplateType enumerates supported template engines.
type TemplateType int

const (
	// SprigTemplateType uses Go text/template with the Sprig function library plus SKA helpers.
	SprigTemplateType TemplateType = iota
	// JinjaTemplateType uses a Jinja-like engine backed by pongo2.
	JinjaTemplateType
)
