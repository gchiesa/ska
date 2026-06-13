package stringprocessor

import (
	"log/slog"

	"github.com/gchiesa/ska/pkg/templateprovider"
)

type StringProcessor struct {
	template templateprovider.TemplateService
	log      *slog.Logger
}
