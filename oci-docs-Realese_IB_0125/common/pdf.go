package common

import (
	"bytes"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergePDF(pdfList [][]byte) ([]byte, error) {
	rs := make([]io.ReadSeeker, len(pdfList))
	for i, item := range pdfList {
		rs[i] = bytes.NewReader(item)
	}
	buf := new(bytes.Buffer)

	if err := api.Merge(rs, buf, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
