package templateprovider

import (
	sprig "github.com/go-task/slim-sprig"
	"github.com/palantir/stacktrace"
	"io"
	"os"
	"strings"
	"text/template"
)

type SprigTemplate struct {
	templateFilePath string
	templateContent  string
	variables        map[string]interface{}
	textTemplate     *template.Template
}

func NewSprigTemplate(name string) *SprigTemplate {
	t := template.New(name)
	t.Funcs(sprig.FuncMap())
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
