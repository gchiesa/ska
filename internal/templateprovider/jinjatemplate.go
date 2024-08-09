package templateprovider

import (
	"github.com/flosch/pongo2/v6"
	"github.com/palantir/stacktrace"
	"io"
)

type JinjaTemplate struct {
	templateFilePath string
	templateContent  string
	variables        map[string]interface{}
	pongo2Template   *pongo2.Template
}

func NewJinjaTemplate(name string) *JinjaTemplate {
	return &JinjaTemplate{}
}

func (t *JinjaTemplate) FromString(templateContent string) error {
	t.templateContent = templateContent
	tpl, err := pongo2.FromString(t.templateContent)
	if err != nil {
		return stacktrace.Propagate(err, "failed to parse template")
	}
	t.pongo2Template = tpl
	return nil
}

func (t *JinjaTemplate) FromFile(templateFilePath string) error {
	tpl, err := pongo2.FromFile(templateFilePath)
	if err != nil {
		return err
	}
	t.pongo2Template = tpl
	return nil
}

func (t *JinjaTemplate) Execute(fp io.Writer, withVariables map[string]interface{}) error {
	t.variables = withVariables
	var context = make(pongo2.Context)
	for k, v := range t.variables {
		context[k] = v.(any)
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

func (t *JinjaTemplate) WithErrorOnMissingKey(_ bool) {
}

func (t *JinjaTemplate) IsMissingKeyError(err error) bool {
	return err.Error() == "TokenError"
}
