package docx

import (
	"bytes"
	"text/template"
)

func RenderDocxTemplate(templatePath string, funcMap template.FuncMap, data interface{}) ([]byte, error) {
	docTmpl, err := GetTemplate(templatePath)
	if err != nil {
		return nil, err
	}
	defer docTmpl.Document.Close()

	if funcMap != nil {
		docTmpl.AddFunctions(funcMap)
	}

	if err := docTmpl.Parse(); err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)

	if err := docTmpl.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
