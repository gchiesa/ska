package contentprovider

import (
	"github.com/apex/log"
	"github.com/otiai10/copy"
	"os"
)

type LocalPath struct {
	sourcePath string
	workingDir string
	log        *log.Entry
}

const LocalPathPrefix = "file://"

func NewLocalPath(path string) (*LocalPath, error) {
	logCtx := log.WithFields(log.Fields{
		"pkg":  "contentprovider",
		"type": "github",
	})
	tmpDir, err := os.MkdirTemp("", workingDirPrefix)
	if err != nil {
		return nil, err
	}
	return &LocalPath{
		sourcePath: path,
		workingDir: tmpDir,
		log:        logCtx,
	}, nil
}

func (cp *LocalPath) RemoteURI() string {
	return cp.sourcePath
}

func (cp *LocalPath) Cleanup() error {
	cp.log.WithFields(log.Fields{"workingDir": cp.workingDir}).Debug("removing working dir.")
	err := os.RemoveAll(cp.workingDir)
	return err
}
func (cp *LocalPath) DownloadContent() error {
	if cp.workingDir == "" {
		var err error
		cp.workingDir, err = os.MkdirTemp("", workingDirPrefix)
		if err != nil {
			return err
		}
	}
	// recursive copy from cp.sourcePath to cp.workingDir
	return copy.Copy(cp.sourcePath, cp.workingDir, copy.Options{PreserveTimes: false})
}

func (cp *LocalPath) WorkingDir() string {
	return cp.workingDir
}
