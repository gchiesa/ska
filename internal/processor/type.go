package processor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/multipart"
	"github.com/gchiesa/ska/internal/templateservice"
)

type FileTreeProcessor struct {
	sourcePath             string
	sourceIgnorePaths      []string
	destinationPathRoot    string
	destinationIgnorePaths []string
	workingDir             string
	multiparts             []*multipart.Multipart
	templateService        *templateservice.SprigTemplate
	log                    *log.Entry
}
