package templateprovider

import (
	"errors"
	"github.com/Masterminds/sprig/v3"
	"github.com/palantir/stacktrace"
	"io"
	"os"
	"strings"
	"text/template"
)

// ErrOptionalContinue is returned by the optional function when a value is
// considered empty and the template should skip emitting it.
var ErrOptionalContinue = errors.New("optional field was empty")

// SprigTemplate implements TemplateService using Go's text/template enriched
// with the Sprig function library plus SKA-specific helpers (optional, empty,
// notempty).
type SprigTemplate struct {
	templateFilePath string
	templateContent  string
	variables        map[string]interface{}
	textTemplate     *template.Template
}

// NewSprigTemplate creates a Sprig-based template engine instance using the
// provided name.
func NewSprigTemplate(name string) *SprigTemplate {
	t := template.New(name)
	t.Funcs(sprig.FuncMap())
	// add specific SKA functions
	skaFunctions := map[string]interface{}{
		"optional": func(c bool, v string) (string, error) {
			if !c {
				return "", ErrOptionalContinue
			}
			if strings.TrimSpace(v) == "" {
				return "", ErrOptionalContinue
			}
			return v, nil
		},
		"empty": func(v string) bool {
			return strings.TrimSpace(v) == ""
		},
		"notempty": func(v string) bool {
			return strings.TrimSpace(v) != ""
		},
	}
	t.Funcs(skaFunctions)
	return &SprigTemplate{
		textTemplate: t,
	}
}

// WithErrorOnMissingKey sets the behavior of the underlying template engine to
// error when a key is missing, if state is true.
func (t *SprigTemplate) WithErrorOnMissingKey(state bool) {
	if state {
		t.textTemplate.Option("missingkey=error")
	} else {
		t.textTemplate.Option("missingkey=default")
	}
}

// FromString parses the given template string.
func (t *SprigTemplate) FromString(templateContent string) error {
	t.templateContent = templateContent
	tpl, err := t.textTemplate.Parse(t.templateContent)
	if err != nil {
		return stacktrace.Propagate(err, "failed to parse template")
	}
	t.textTemplate = tpl
	return nil
}

// FromFile loads and parses a template from a file path.
func (t *SprigTemplate) FromFile(templateFilePath string) error {
	fileContent, err := os.ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	t.templateFilePath = templateFilePath
	t.templateContent = string(fileContent)
	tpl, err := t.textTemplate.Parse(t.templateContent)
	if err != nil {
		return stacktrace.Propagate(err, "failed to parse template")
	}
	t.textTemplate = tpl
	return nil
}

// Execute renders the template into the writer using the provided variables.
func (t *SprigTemplate) Execute(fp io.Writer, withVariables map[string]interface{}) error {
	t.variables = withVariables
	return t.textTemplate.Execute(fp, t.variables)
}

// IsMissingKeyError reports whether the error is due to a missing key.
func (t *SprigTemplate) IsMissingKeyError(err error) bool {
	return strings.Contains(err.Error(), "map has no entry for key")
}

// IsOptionalError reports whether the error corresponds to an optional skip.
func (t *SprigTemplate) IsOptionalError(err error) bool {
	return errors.Is(err, ErrOptionalContinue)
}
