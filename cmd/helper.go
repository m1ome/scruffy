package cmd

import "os"

type WorkingDirGetter func() (string, error)

func Getwd() (string, error) {
	return os.Getwd()
}
