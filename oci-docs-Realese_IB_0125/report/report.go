package report

import "io"

type Report interface {
	Render() error
	WriteTo(w io.Writer) (int64, error)
}
