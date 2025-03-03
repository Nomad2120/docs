package model

// PDFSettings -
type PDFSettings struct {
	Dpi              *uint    `json:"dpi,omitempty" example:"300"`                // Change the dpi explicitly (this has no effect on X11 based systems)
	Grayscale        *bool    `json:"grayscale,omitempty" example:"false"`        // PDF will be generated in grayscale
	ImageDpi         *uint    `json:"imageDpi,omitempty" example:"600"`           // When embedding images scale them down to this dpi (default 600)
	ImageQuality     *uint    `json:"imageQuality,omitempty" example:"94"`        // When jpeg compressing images use this quality (default 94)
	MarginBottom     *uint    `json:"marginBottom,omitempty"  example:"10"`       // Set the page bottom margin
	MarginLeft       *uint    `json:"marginLeft,omitempty"  example:"10"`         // Set the page left margin (default 10mm)
	MarginRight      *uint    `json:"marginRight,omitempty"  example:"10"`        // Set the page right margin (default 10mm)
	MarginTop        *uint    `json:"marginTop,omitempty"  example:"10"`          // Set the page top margin
	NoCollate        *bool    `json:"noCollate,omitempty" example:"false"`        // Do not collate when printing multiple copies (default collate)
	NoPdfCompression *bool    `json:"noPdfCompression,omitempty" example:"false"` // Do not use lossless compression on pdf objects
	Orientation      *string  `json:"orientation,omitempty" example:"Portrait"`   // Set orientation to Landscape or Portrait (default Portrait)
	PageHeight       *uint    `json:"pageHeight,omitempty"`                       // Page height
	PageSize         *string  `json:"pageSize,omitempty" example:"A4"`            // Set paper size to: A4, Letter, etc. (default A4)
	PageWidth        *uint    `json:"pageWidth,omitempty"`                        // Page width
	Encoding         *string  `json:"encoding,omitempty"  example:"utf-8"`        // Set the default text encoding, for input
	MinimumFontSize  *uint    `json:"minimumFontSize,omitempty"`                  // Minimum font size
	FooterFontSize   *uint    `json:"footerFontSize,omitempty" example:"12"`      // Set footer font size (default 12)
	Zoom             *float64 `json:"zoom,omitempty" example:"1"`                 // Use this zoom factor (default 1)
}

type Html2PdfRequest struct {
	Settings   PDFSettings `json:"settings"`
	HTMLBase64 string      `json:"htmlBase64"`
}
