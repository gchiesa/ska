package stringprocessor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/templateprovider"
)

type StringProcessor struct {
	template templateprovider.TemplateService
	log      *log.Entry
}
