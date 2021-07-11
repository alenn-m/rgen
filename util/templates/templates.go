package templates

import (
	"bytes"
	"go/format"
	"text/template"
)

// ParseTemplate parses template
func ParseTemplate(tmpl string, data interface{}, functions map[string]interface{}) (string, error) {
	modelTemplate, err := template.New("template").
		Funcs(functions).
		Parse(tmpl)
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	err = modelTemplate.Execute(&buff, data)
	if err != nil {
		return "", err
	}

	content, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	return string(content), nil
}
