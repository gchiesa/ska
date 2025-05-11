package templateprovider

import (
	"io"
)

type TemplateService interface {
	// FromFile Load template from file
	FromFile(path string) error
	// FromString Load template from string
	FromString(templateContent string) error
	// Execute Execute the template
	Execute(fp io.Writer, withVariables map[string]interface{}) error
	// WithErrorOnMissingKey Set error on missing key
	WithErrorOnMissingKey(key bool)
	IsMissingKeyError(err error) bool
	IsOptionalError(err error) bool
}

type TemplateType int

const (
	SprigTemplateType TemplateType = iota
	JinjaTemplateType
)
