package templateservice

import (
	"io"
)

type TemplateService interface {
	FromFile(path string) error
	FromString(templateContent string) error
	Execute(fp io.Writer, withVariables map[string]interface{}) error
}
