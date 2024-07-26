package processor

import (
	"bytes"
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/multipart"
	"github.com/otiai10/copy"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// buildStagingFileTree build the staging file tree
// allocate the folders and files by copying a local source (upstream) blueprint into a destination
// so that they are ready to be rendered
// if file structure (directories, files) are templated then they are expanded and rendered
func (tp *FileTreeProcessor) buildStagingFileTree(withVariables map[string]interface{}) error {
	logger := tp.log.WithFields(log.Fields{"method": "buildStagingFileTree"})

	// walk the sourcePath and render the files
	sPathAbs, err := filepath.Abs(tp.sourcePath)
	if err != nil {
		return err
	}
	dPathAbs, err := filepath.Abs(tp.workingDir)
	if err != nil {
		return err
	}

	err = filepath.Walk(tp.sourcePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		sRelPath, err := filepath.Rel(sPathAbs, path)
		if err != nil {
			return err
		}

		// filter out path that should not be processed
		if !tp.shouldProcessFile(sRelPath) {
			logger.WithFields(log.Fields{"path": sRelPath}).Debug("skipping path")
			return nil
		}

		// create a template from the file name as it was a template
		if err := tp.templateService.FromString(sRelPath); err != nil {
			return err
		}

		// render the template
		buff := bytes.NewBufferString("")
		if err := tp.templateService.Execute(buff, withVariables); err != nil {
			if tp.templateService.IsMissingKeyError(err) {
				logger.WithFields(log.Fields{"path": sRelPath}).Errorf("missing variable while rendering file path: %s", sRelPath)
			}
			return err
		}

		dPath := filepath.Join(dPathAbs, buff.String())

		// if it's file we copy the file to the destination
		if !info.IsDir() {
			if err := copy.Copy(path, dPath, copy.Options{PreserveTimes: false}); err != nil {
				return err
			}
		} else {
			// if directory we allocate all the path
			if err := os.MkdirAll(dPath, 0o755); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// loadMultiparts decompose the files in the staging directory that are managed partials
// create a set of partials that are related to the files in the staging directory
func (tp *FileTreeProcessor) loadMultiparts() error {
	logger := tp.log.WithFields(log.Fields{"method": "loadMultiparts"})
	err := filepath.Walk(tp.workingDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(tp.workingDir, path)
		if err != nil {
			return err
		}

		if !tp.shouldProcessFile(relPath) {
			tp.log.WithFields(log.Fields{"path": relPath}).Debug("skipping path")
			return nil
		}

		// if it's file we copy the file to the destination
		if !info.IsDir() {
			multipart, err := multipart.NewMultipartFromFile(absPath, relPath)
			if err != nil {
				return err
			}

			if err = multipart.ParseParts(); err != nil { //nolint:gocritic
				return err
			}
			files, err := multipart.PartsToFiles()
			if err != nil {
				return err
			}
			logger.WithFields(log.Fields{"parts": files, "multipart": relPath}).Debug("Generating Parts from Multipart.")
			tp.multiparts = append(tp.multiparts, multipart)
		} else {
			logger.WithFields(log.Fields{"filePath": relPath}).Debug("Skipping because is a directory.")
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// renderStagingFileTree render all the templates, but if a partial exists for a file then expands only the partials
func (tp *FileTreeProcessor) renderStagingFileTree(withVariables map[string]interface{}) error {
	logger := tp.log.WithFields(log.Fields{"method": "renderStagingFileTree"})
	err := filepath.Walk(tp.workingDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(tp.workingDir, path)
		if err != nil {
			return err
		}

		// if it's file we copy the file to the destination
		if !info.IsDir() {
			// check if the file is to process
			if tp.multipartExistsAndHasPartials(relPath) {
				logger.WithFields(log.Fields{"filePath": relPath}).Debug("Skipping file because it's a Multipart with Parts.")
				return nil
			}

			if err := tp.templateService.FromFile(absPath); err != nil {
				return err
			}

			// render the template
			buff := bytes.NewBufferString("")
			if err := tp.templateService.Execute(buff, withVariables); err != nil {
				if tp.templateService.IsMissingKeyError(err) {
					logger.WithFields(log.Fields{"path": relPath}).Errorf("missing variable while rendering file: %s", relPath)
				}
				return err
			}

			logger.WithFields(log.Fields{"filePath": relPath}).Debug("Saving rendered file.")
			if err := os.WriteFile(absPath, []byte(buff.String()), 0o644); err != nil { //nolint:gosimple // we don't need to check the error here
				return err
			}
		} else {
			logger.WithFields(log.Fields{"filePath": relPath}).Debug("Skipping because is a directory directory.")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (tp *FileTreeProcessor) shouldProcessFile(path string) bool {
	if strings.HasPrefix(path, ".git/") {
		return false
	}
	fileParts := strings.Split(path, configuration.AppIdentifier)

	// if it's a file in the form of `file.swanson-....` we skip it
	if len(fileParts) > 1 && strings.HasSuffix(fileParts[0], ".") && strings.HasPrefix(fileParts[1], "-") {
		return false
	}
	return true
}

func (tp *FileTreeProcessor) getMultipartByID(id string) (*multipart.Multipart, error) {
	for _, pc := range tp.multiparts {
		if pc.ID() == id {
			return pc, nil
		}
	}
	return nil, fmt.Errorf("multipart not found: %s", id)
}

func (tp *FileTreeProcessor) multipartWithIDExists(id string) bool {
	if _, err := tp.getMultipartByID(id); err != nil {
		return false
	}
	return true
}

func (tp *FileTreeProcessor) multipartExistsAndHasPartials(id string) bool {
	if !tp.multipartWithIDExists(id) {
		return false
	}
	pc, _ := tp.getMultipartByID(id)
	return pc.HasParts()
}

func (tp *FileTreeProcessor) fileIsMultipart(relativeFilePath string) bool {
	return tp.multipartWithIDExists(relativeFilePath)
}

// updateDestinationFileTree assemble partials and update the destination directory
// what it does:
// copy the staging directory to the destination with the process below:
// - for each file (non-partial) copy the file first
// - then replace the partials with the expanded content
// - copy to destination with the logic:
//   - if the file mustBeSkipped then skip, otherwise copy
//   - if the file already exists in the destination then
//     only replace the partials with the expanded content
func (tp *FileTreeProcessor) updateDestinationFileTree() error {
	logger := tp.log.WithFields(log.Fields{"method": "updateDestinationFileTree"})
	err := filepath.Walk(tp.workingDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(tp.workingDir, path)
		if err != nil {
			return err
		}

		// if it's file we copy the file to the destination
		if !info.IsDir() {
			// is it is not a swanson file we copy to destination
			if !tp.shouldProcessFile(relPath) {
				logger.WithFields(log.Fields{"filePath": relPath}).Debug("Skipping file because should not be processed.")
				return nil
			}
			// not a managed file then we copy it to destination
			if !tp.fileIsMultipart(relPath) {
				logger.WithFields(log.Fields{"filePath": relPath, "destination": tp.destinationPathRoot}).Debug("Copying file to destination.")
				if err := copy.Copy(absPath, filepath.Join(tp.destinationPathRoot, relPath)); err != nil {
					return err
				}
				return nil
			}
			mp, _ := tp.getMultipartByID(relPath)

			// if it has no partials then we just copy as normal expanded file
			if !mp.HasParts() {
				logger.WithFields(log.Fields{"filePath": relPath, "destination": tp.destinationPathRoot}).Debug("Copying non multipart file to destination.")
				if err := copy.Copy(absPath, filepath.Join(tp.destinationPathRoot, relPath)); err != nil {
					return err
				}
				return nil
			}

			// assemble back the partial container with the rendered partials
			logger.WithFields(log.Fields{"filePath": relPath, "destination": tp.destinationPathRoot}).Debug("Compiling Multipart file to destination.")
			if err := mp.CompileToFile(filepath.Join(tp.destinationPathRoot, relPath), false); err != nil {
				return err
			}
		} else {
			// if directory we allocate all the path
			if err := os.MkdirAll(filepath.Join(tp.destinationPathRoot, relPath), 0o755); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
