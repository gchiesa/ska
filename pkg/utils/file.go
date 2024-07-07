package utils

import "os"

func CreateWorkingDir(dir, pattern string) (string, error) {
	dir, err := os.MkdirTemp("", "swansonRenderer")
	if err != nil {
		return "", err
	}
	return dir, nil
}
