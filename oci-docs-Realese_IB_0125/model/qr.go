package model

// FillQRPageRequest -
type FillQRPageRequest struct {
	ID int `json:"id" uri:"id" binding:"required"`
}

// FillActResponse -
type FillQRPageResponse struct {
	BaseResponse
	Result *FillQRPageResult `json:"result"`
}

// FillQRPageResult -
type FillQRPageResult struct {
	PDFBase64 string `json:"pdfBase64"`
}

type QRPage struct {
	Name     string
	Address  string
	QRBase64 string
}
