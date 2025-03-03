package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"text/template"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/skip2/go-qrcode"

	"github.com/pkg/errors"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	store "gitlab.enterprise.qazafn.kz/oci/oci-docs/template"
)

func (c *appService) FillQRPage(ctx context.Context, req *model.FillQRPageRequest) (*model.FillQRPageResult, error) {
	action := "FillQRPage"
	osi, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}
	t := template.New(action)
	t, err = t.Parse(string(store.GetQrPageTemplate()))
	if err != nil {
		return nil, errors.Wrap(err, action)
	}
	b, err := qrcode.Encode(fmt.Sprintf("%s%d", "https://t.me/OSISubscriptionBot?start=", req.ID), qrcode.Medium, 256)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	qrPage := model.QRPage{
		Name:     osi.Name,
		Address:  osi.Address,
		QRBase64: base64.StdEncoding.EncodeToString(b),
	}

	var buf bytes.Buffer
	t.Execute(&buf, &qrPage)
	html := base64.StdEncoding.EncodeToString(buf.Bytes())

	no := false
	var margin uint = 5
	zoom := 1.4
	portret := string(pdf.OrientationPortrait)
	settings := model.PDFSettings{
		Grayscale:    &no,
		MarginLeft:   &margin,
		MarginTop:    &margin,
		MarginBottom: &margin,
		Zoom:         &zoom,
		Orientation:  &portret,
	}
	pdfB64, err := c.doc.HTMLToPDFBase64(html, &settings)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	return &model.FillQRPageResult{
		PDFBase64: pdfB64,
	}, nil
}
