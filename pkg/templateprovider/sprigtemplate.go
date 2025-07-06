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

var ErrOptionalContinue = errors.New("optional field was empty")

type SprigTemplate struct {
	templateFilePath string
	templateContent  string
	variables        map[string]interface{}
	textTemplate     *template.Template
}

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

func (t *SprigTemplate) WithErrorOnMissingKey(state bool) {
	if state {
		t.textTemplate.Option("missingkey=error")
	} else {
		t.textTemplate.Option("missingkey=default")
	}
}

func (t *SprigTemplate) FromString(templateContent string) error {
	t.templateContent = templateContent
	tpl, err := t.textTemplate.Parse(t.templateContent)
	if err != nil {
		return stacktrace.Propagate(err, "failed to parse template")
	}
	t.textTemplate = tpl
	return nil
}

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

func (t *SprigTemplate) Execute(fp io.Writer, withVariables map[string]interface{}) error {
	t.variables = withVariables
	return t.textTemplate.Execute(fp, t.variables)
}

func (t *SprigTemplate) IsMissingKeyError(err error) bool {
	return strings.Contains(err.Error(), "map has no entry for key")
}

func (t *SprigTemplate) IsOptionalError(err error) bool {
	return errors.Is(err, ErrOptionalContinue)
}
