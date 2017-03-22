package cmd

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const errEmptyResult = "Parsed empty template"
const errEmptyDirectory = "Empty input directory"
const defaultTemplate = "index"

// Parser empty parser object
type Parser struct {
	Wd WorkingDirGetter
}

// NewParser returns new Parser object
func NewParser() *Parser {
	return &Parser{
		Wd: Getwd,
	}
}

// Parse parsing folder with env variables passed in
// by default parsing will be done though template/text
// Golang package.
func (p Parser) Parse(input string, env interface{}) (res []byte, err error) {
	if !filepath.IsAbs(input) {
		cwd, cwdErr := p.Wd()
		if cwdErr != nil {
			err = cwdErr
			return
		}

		input = path.Join(cwd, input)
	}

	// Traversing path
	files := []string{}
	err = filepath.Walk(input, func(path string, info os.FileInfo, fileError error) error {
		if fileError != nil {
			return fileError
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return
	}

	if len(files) == 0 {
		err = errors.New(errEmptyDirectory)
		return
	}

	// Building templates
	t, err := template.ParseFiles(files...)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, defaultTemplate, env)
	if err != nil {
		return
	}

	res = buf.Bytes()
	if len(res) == 0 {
		err = errors.New(errEmptyResult)
		return
	}

	return
}
