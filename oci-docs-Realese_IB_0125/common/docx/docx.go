package docx

import (
	"archive/zip"
	"errors"
	"io"
)

// Docx struct that contains data from a docx
type Docx struct {
	zipReader *zip.ReadCloser
	content   string
}

// ReadFile func reads a docx file
func (d *Docx) ReadFile(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return errors.New("cannot Open File")
	}
	content, err := readText(reader.File)
	if err != nil {
		return errors.New("cannot Read File")
	}
	d.zipReader = reader
	if content == "" {
		return errors.New("file has no content")
	}
	d.content = cleanText(content)
	return nil
}

// UpdateContent updates the content string
func (d *Docx) UpdateContent(newContent string) {
	d.content = newContent
}

// GetContent returns the string content
func (d *Docx) GetContent() string {
	return d.content
}

func (d *Docx) Write(ioWriter io.Writer, data string) error {
	// reformat string, for some reason the first char is converted to &lt;
	w := zip.NewWriter(ioWriter)
	for _, file := range d.zipReader.File {
		var writer io.Writer
		var readCloser io.ReadCloser
		writer, err := w.Create(file.Name)
		if err != nil {
			return err
		}
		readCloser, err = file.Open()
		if err != nil {
			return err
		}
		if file.Name == "word/document.xml" {
			_, err := writer.Write([]byte(data))
			if err != nil {
				return err
			}
		} else {
			b, err := streamToByte(readCloser)
			if err != nil {
				return err
			}
			_, err = writer.Write(b)
			if err != nil {
				return err
			}
		}
	}
	return w.Close()
}

// Close the document
func (d *Docx) Close() error {
	return d.zipReader.Close()
}
