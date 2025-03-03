package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
)

// readText reads text from a word document
func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = retrieveWordDoc(files)
	if err != nil {
		return text, err
	}
	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}

	text, err = wordDocToString(documentReader)
	return
}

// wordDocToString converts a word document to string
func wordDocToString(reader io.Reader) (string, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// retrieveWordDoc fetches a word document.
func retrieveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml file not found")
	}
	return
}

func streamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(stream); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// normalize fixes quotation marks in documnet
func normalizeQuotes(in rune) rune {
	switch in {
	case '“', '”':
		return '"'
	case '‘', '’':
		return '\''
	}
	return in
}

// cleans template tagged text of all brakets
func normalizeAll(text string) string {
	brakets := regexp.MustCompile("<.*?>")
	quotes := regexp.MustCompile("&quot;")
	text = brakets.ReplaceAllString(text, "")
	text = quotes.ReplaceAllString(text, "\"")
	return strings.Map(normalizeQuotes, text)
}

func cleanText(text string) string {
	braketFinder := regexp.MustCompile("{{.*?}}")
	return braketFinder.ReplaceAllStringFunc(text, normalizeAll)
}

func escapeData(t interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovery %v", r)
		}
	}()
	v := reflect.ValueOf(t)
	if v.Type().Kind() != reflect.Ptr {
		err = errors.New("need pass pointer of struct")
		return
	}

	for {
		if v.Type().Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}

	err = escapeStruct(v)
	return
}

func escapeStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		if err := escapeField(f); err != nil {
			return err
		}
	}
	return nil
}

func escapeSlice(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		if err := escapeField(v.Index(i)); err != nil {
			return nil
		}
	}
	return nil
}

func escapeMap(v reflect.Value) error {
	for _, k := range v.MapKeys() {

		if v.MapIndex(k).Kind() == reflect.String {
			b := new(bytes.Buffer)
			if err := xml.EscapeText(b, []byte(v.MapIndex(k).String())); err != nil {
				return err
			}
			v.SetMapIndex(k, reflect.ValueOf(b.String()))
		} else {
			if err := escapeField(v.MapIndex(k)); err != nil {
				return nil
			}
			v.SetMapIndex(k, v.MapIndex(k))
		}

	}
	return nil
}

func escapeField(f reflect.Value) error {
	if f.Type().Kind() == reflect.Ptr {
		f = f.Elem()
	}

	switch f.Kind() {
	case reflect.Struct:
		if err := escapeStruct(f); err != nil {
			return nil
		}
	case reflect.Slice:
		if err := escapeSlice(f); err != nil {
			return nil
		}
	case reflect.Map:
		if err := escapeMap(f); err != nil {
			return nil
		}
	case reflect.String:
		b := new(bytes.Buffer)
		if err := xml.EscapeText(b, []byte(f.String())); err != nil {
			return err
		}
		f.SetString(b.String())
	default:
		return nil
	}

	return nil
}
