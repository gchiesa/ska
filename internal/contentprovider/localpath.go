package contentprovider

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/otiai10/copy"
)

type LocalPath struct {
	sourcePath string
	workingDir string
	log        *slog.Logger
}

const LocalPathPrefix = "file://"

func NewLocalPath(path string, logger *slog.Logger) (*LocalPath, error) {
	if logger == nil {
		logger = slog.Default()
	}
	logCtx := logger.With(logFieldPkg, logPkg, logFieldType, "localpath")
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
	return fmt.Sprintf("%s%s", LocalPathPrefix, cp.sourcePath)
}

func (cp *LocalPath) Cleanup() error {
	cp.log.With(logFieldWorkingDir, cp.workingDir).Debug("removing working dir.")
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
