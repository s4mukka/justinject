package utils

import (
	"bytes"
	"html/template"
)

var (
	templateParseFiles = template.ParseFiles
)

func (u *Utils) TemplateToString(path string, data any) (string, error) {
	tmpl, err := templateParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
