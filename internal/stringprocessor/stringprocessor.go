package stringprocessor

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/gchiesa/ska/pkg/templateprovider"
)

func NewStringProcessor(options ...func(stringProcessor *StringProcessor)) *StringProcessor {
	sp := &StringProcessor{
		template: nil,
		log:      slog.Default().With("pkg", "string-processor"),
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

// WithLogger injects a *slog.Logger into the processor.
// The processor will add its own "pkg" field to the received logger.
func WithLogger(logger *slog.Logger) func(sp *StringProcessor) {
	return func(sp *StringProcessor) {
		if logger != nil {
			sp.log = logger.With("pkg", "string-processor")
		}
	}
}

func (sp *StringProcessor) Render(text string, withVariables map[string]interface{}) (string, error) {
	logger := sp.log.With("method", "Render")

	if err := sp.template.FromString(text); err != nil {
		return "", err
	}

	// render the template
	buff := bytes.NewBufferString("")
	if err := sp.template.Execute(buff, withVariables); err != nil {
		if sp.template.IsMissingKeyError(err) {
			logger.With("error", err.Error()).Error(fmt.Sprintf("missing variable while rendering string: %v", err))
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
