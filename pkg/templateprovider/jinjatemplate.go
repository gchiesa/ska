package templateprovider

import (
	"io"

	"github.com/flosch/pongo2/v6"
	"github.com/palantir/stacktrace"
)

// JinjaTemplate implements TemplateService using the pongo2 Jinja-like engine.

type JinjaTemplate struct {
	templateContent string
	variables       map[string]interface{}
	pongo2Template  *pongo2.Template
}

// NewJinjaTemplate creates a Jinja-like template engine instance backed by pongo2.
func NewJinjaTemplate(_ string) *JinjaTemplate {
	return &JinjaTemplate{}
}

// FromString parses a Jinja-like template from the provided string.
func (t *JinjaTemplate) FromString(templateContent string) error {
	t.templateContent = templateContent
	tpl, err := pongo2.FromString(t.templateContent)
	if err != nil {
		return stacktrace.Propagate(err, "failed to parse template")
	}
	t.pongo2Template = tpl
	return nil
}

// FromFile parses a Jinja-like template from a file.
func (t *JinjaTemplate) FromFile(templateFilePath string) error {
	tpl, err := pongo2.FromFile(templateFilePath)
	if err != nil {
		return err
	}
	t.pongo2Template = tpl
	return nil
}

// Execute renders the Jinja-like template into the writer using variables.
func (t *JinjaTemplate) Execute(fp io.Writer, withVariables map[string]interface{}) error {
	t.variables = withVariables
	var context = make(pongo2.Context)
	for k, v := range t.variables {
		context[k] = v
	}
	renderedContent, err := t.pongo2Template.Execute(context)
	if err != nil {
		return err
	}

	_, err = fp.Write([]byte(renderedContent))
	if err != nil {
		return err
	}
	return nil
}

// WithErrorOnMissingKey is a no-op for JinjaTemplate as pongo2 handles
// missing keys according to its own semantics.
func (t *JinjaTemplate) WithErrorOnMissingKey(_ bool) {
}

// IsMissingKeyError reports whether the error indicates a missing key for
// pongo2 execution.
func (t *JinjaTemplate) IsMissingKeyError(err error) bool {
	return err.Error() == "TokenError"
}

// IsOptionalError always returns false for JinjaTemplate; optional skip is not
// signaled via errors in this engine.
func (t *JinjaTemplate) IsOptionalError(_ error) bool {
	return false
}
