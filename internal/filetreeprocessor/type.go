package filetreeprocessor

import (
	"log/slog"

	"github.com/gchiesa/ska/internal/multipart"
	"github.com/gchiesa/ska/pkg/templateprovider"
)

const (
	logFieldMethod      = "method"
	logFieldPath        = "path"
	logFieldFilePath    = "filePath"
	logFieldDestination = "destination"
)

type FileTreeProcessor struct {
	sourcePath             string
	sourceIgnorePaths      []string
	destinationPathRoot    string
	destinationIgnorePaths []string
	workingDir             string
	multiparts             []*multipart.Multipart
	template               templateprovider.TemplateService
	log                    *slog.Logger
}
