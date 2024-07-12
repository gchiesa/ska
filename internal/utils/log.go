package utils

import "github.com/apex/log"

func GetLoggerFor() *log.Entry {
	return log.WithFields(log.Fields{
		"pkg": "utils",
	})
}
