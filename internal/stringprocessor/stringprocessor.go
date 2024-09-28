package stringprocessor

import (
	"bytes"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/templateprovider"
)

func NewStringProcessor(options ...func(stringProcessor *StringProcessor)) *StringProcessor {
	logCtx := log.WithFields(log.Fields{
		"pkg": "string-processor",
	})

	sp := &StringProcessor{
		template: nil,
		log:      logCtx,
	}
	// configure options
	for _, opt := range options {
		opt(sp)
	}
	return sp
}

func WithTemplateService(ts templateprovider.TemplateService) func(sp *StringProcessor) {
	return func(sp *StringProcessor) {
		sp.template = ts
	}
}

func (sp *StringProcessor) Render(text string, withVariables map[string]interface{}) (string, error) {
	logger := sp.log.WithFields(log.Fields{"method": "Render"})

	if err := sp.template.FromString(text); err != nil {
		return "", err
	}

	// render the template
	buff := bytes.NewBufferString("")
	if err := sp.template.Execute(buff, withVariables); err != nil {
		if sp.template.IsMissingKeyError(err) {
			logger.WithFields(log.Fields{"error": err.Error()}).Errorf("missing variable while rendering string: %v", err)
		}
		return "", err
	}
	return buff.String(), nil
}

func (sp *StringProcessor) RenderSliceOfStrings(text []string, variables map[string]interface{}) ([]string, error) {
	var result []string
	for _, entry := range text {
		renderedEntry, err := sp.Render(entry, variables)
		if err != nil {
			return nil, err
		}
		result = append(result, renderedEntry)
	}
	return result, nil
}
