package types

import (
	"bytes"
	"html/template"
)

type Replacements struct {
	Name string
}

func get_replacements(config Config) Replacements {
	return Replacements{
		Name: config.Project_name,
	}
}

func (r Replacements) in_string(target string) (string, error) {
	var buffer bytes.Buffer
	err := template.Must(template.New("").Parse(target)).Execute(&buffer, r)
	return buffer.String(), err
}
