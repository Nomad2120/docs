package conv

import (
	"bytes"
	"encoding/base64"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

// DocConverter -
type DocConverter interface {
	HTMLToPDFBase64(htmlDocBase64 string, settings *model.PDFSettings) (string, error)
	HTMLToPDF(htmlBytes []byte, settings *model.PDFSettings) ([]byte, error)
}

type wkDocConverter struct {
}

func NewWKDocConverter() (DocConverter, error) {
	return &wkDocConverter{}, nil
}

// func (w *wkDocConverter) HTMLToPDFBase64(htmlDocBase64 string) (string, error) {
// 	// Create new PDF generator
// 	pdfg, err := pdf.NewPDFGenerator()
// 	if err != nil {
// 		return "", err
// 	}

// 	// Set global options
// 	pdfg.Dpi.Set(300)
// 	pdfg.Orientation.Set(pdf.OrientationLandscape)
// 	pdfg.Grayscale.Set(false)

// 	b, err := base64.StdEncoding.DecodeString(htmlDocBase64)
// 	if err != nil {
// 		return "", err
// 	}
// 	page := pdf.NewPageReader(strings.NewReader(string(b)))

// 	// Set options for this page
// 	page.FooterRight.Set("[page]")
// 	page.FooterFontSize.Set(16)
// 	page.Zoom.Set(1.5)

// 	// Add to document
// 	pdfg.AddPage(page)

// 	// Create PDF document in internal buffer
// 	if err := pdfg.Create(); err != nil {
// 		return "", err
// 	}

// 	encoded := base64.StdEncoding.EncodeToString(pdfg.Buffer().Bytes())
// 	return encoded, nil
// }

func (w *wkDocConverter) HTMLToPDFBase64(htmlDocBase64 string, settings *model.PDFSettings) (string, error) {
	htmlBytes, err := base64.StdEncoding.DecodeString(htmlDocBase64)
	if err != nil {
		return "", err
	}
	b, err := w.HTMLToPDF(htmlBytes, settings)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded, nil
}

func (w *wkDocConverter) HTMLToPDF(htmlBytes []byte, settings *model.PDFSettings) ([]byte, error) {
	// Create new PDF generator
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(pdf.OrientationPortrait)
	pdfg.Grayscale.Set(true)

	if settings != nil {
		if settings.Dpi != nil {
			pdfg.Dpi.Set(*settings.Dpi)
		}
		if settings.Grayscale != nil {
			pdfg.Grayscale.Set(*settings.Grayscale)
		}
		if settings.ImageDpi != nil {
			pdfg.ImageDpi.Set(*settings.ImageDpi)
		}
		if settings.ImageQuality != nil {
			pdfg.ImageQuality.Set(*settings.ImageQuality)
		}
		if settings.MarginBottom != nil {
			pdfg.MarginBottom.Set(*settings.MarginBottom)
		}
		if settings.MarginLeft != nil {
			pdfg.MarginLeft.Set(*settings.MarginLeft)
		}
		if settings.MarginRight != nil {
			pdfg.MarginRight.Set(*settings.MarginRight)
		}
		if settings.MarginTop != nil {
			pdfg.MarginTop.Set(*settings.MarginTop)
		}
		if settings.NoCollate != nil {
			pdfg.NoCollate.Set(*settings.NoCollate)
		}
		if settings.NoPdfCompression != nil {
			pdfg.NoPdfCompression.Set(*settings.NoPdfCompression)
		}
		if settings.Orientation != nil {
			pdfg.Orientation.Set(*settings.Orientation)
		}
		if settings.PageHeight != nil {
			pdfg.PageHeight.Set(*settings.PageHeight)
		}
		if settings.PageSize != nil {
			pdfg.PageSize.Set(*settings.PageSize)
		}
		if settings.PageWidth != nil {
			pdfg.PageWidth.Set(*settings.PageWidth)
		}

	}

	page := pdf.NewPageReader(bytes.NewReader(htmlBytes))

	if settings == nil {
		// Set options for this page
		page.FooterRight.Set("[page]")
		page.FooterFontSize.Set(16)
		page.Zoom.Set(1.0)
	} else {
		if settings.Encoding != nil {
			page.Encoding.Set(*settings.Encoding)
		}
		if settings.MinimumFontSize != nil {
			page.MinimumFontSize.Set(*settings.MinimumFontSize)
		}
		if settings.FooterFontSize != nil {
			page.FooterFontSize.Set(*settings.FooterFontSize)
		}
		if settings.Zoom != nil {
			page.Zoom.Set(*settings.Zoom)
		}
	}

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Buffer().Bytes(), nil
}
