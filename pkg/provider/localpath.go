package provider

import (
	"github.com/otiai10/copy"
	"os"
)

type LocalPath struct {
	path       string
	workingDir string
}

func NewLocalPath(path string, inMemoryFs bool) *LocalPath {
	return &LocalPath{
		path: path,
	}
}

func (p *LocalPath) Path() string {
	return p.path
}

func (p *LocalPath) RemoveWorkingDir() error {
	err := os.RemoveAll(p.workingDir)
	return err
}
func (p *LocalPath) DownloadContent() error {
	if p.workingDir == "" {
		var err error
		p.workingDir, err = os.MkdirTemp("", "swansonDownloader")
		if err != nil {
			return err
		}
	}
	// recursive copy from p.path to p.workingDir
	return copy.Copy(p.path, p.workingDir, copy.Options{PreserveTimes: false})
}

func (p *LocalPath) RemotePath() string {
	return p.path
}

func (p *LocalPath) WorkingDir() string {
	return p.workingDir
}
