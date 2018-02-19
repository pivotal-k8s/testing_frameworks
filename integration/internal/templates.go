package internal

import (
	"bytes"
	"html/template"
)

func RenderTemplates(templates []string, vars interface{}) ([]string, error) {
	var renderedTemplates []string

	for _, arg := range templates {
		t, err := template.New(arg).Parse(arg)
		if err != nil {
			return nil, err
		}

		buf := &bytes.Buffer{}
		err = t.Execute(buf, vars)
		if err != nil {
			return nil, err
		}
		renderedTemplates = append(renderedTemplates, buf.String())
	}

	return renderedTemplates, nil
}
