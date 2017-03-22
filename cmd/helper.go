package cmd

import "os"

// WorkingDirGetter interface for Getwd
type WorkingDirGetter func() (string, error)

// Getwd get current working dir piping os.Getwd for a testing purposes
func Getwd() (string, error) {
	return os.Getwd()
}
