package docx

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"text/template"
)

// Document interface is a combintation of methods use for generic data files
type Document interface {
	ReadFile(string) error
	UpdateContent(string)
	GetContent() string
	Write(ioWriter io.Writer, data string) error
	Close() error
}

// DocTemplate struct combines data and methods from both the Document interface
// and golang's templating library
type DocTemplate struct {
	Template *template.Template
	Document Document
}

// GetTemplate uses the file extension to determine the correct document struct to use
func GetTemplate(filePath string) (*DocTemplate, error) {
	var document Document
	switch filepath.Ext(filePath) {
	case ".docx":
		document = new(Docx)
	default:
		return nil, errors.New("unsupported document type")
	}
	err := document.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return &DocTemplate{Document: document, Template: template.New("docxTemplate")}, nil
}

func (docTemplate *DocTemplate) Execute(wr io.Writer, data interface{}) error {
	if err := escapeData(data); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err := docTemplate.Template.Execute(buf, data)
	if err != nil {
		return err
	}
	return docTemplate.Document.Write(wr, buf.String())
}

// AddFunctions adds functions to the template
func (docTemplate *DocTemplate) AddFunctions(funcMap template.FuncMap) {
	docTemplate.Template = docTemplate.Template.Funcs(funcMap)
}

// Parse parses the template
func (docTemplate *DocTemplate) Parse() error {
	temp, err := docTemplate.Template.Parse(docTemplate.Document.GetContent())
	if err != nil {
		return err
	}
	docTemplate.Template = temp
	return nil
}
