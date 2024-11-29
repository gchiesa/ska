package filetreeprocessor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/multipart"
	"github.com/gchiesa/ska/pkg/templateprovider"
)

type FileTreeProcessor struct {
	sourcePath             string
	sourceIgnorePaths      []string
	destinationPathRoot    string
	destinationIgnorePaths []string
	workingDir             string
	multiparts             []*multipart.Multipart
	template               templateprovider.TemplateService
	log                    *log.Entry
}
