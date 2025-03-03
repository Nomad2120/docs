package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"sort"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
	store "gitlab.enterprise.qazafn.kz/oci/oci-docs/template"
	"gitlab.enterprise.qazafn.kz/oci/oci-docsreport/excel"
)

func (c *appService) GetOSVReport(ctx context.Context, req *model.OSVReportRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	osvResp, err := c.coreRepo.GetOSVByPeriod(ctx, req)
	if err != nil {
		return nil, err
	}

	serviceGrops, err := c.coreRepo.GetServiceGroups(ctx)
	if err != nil {
		return nil, err
	}

	rep := excel.NewSaldoXLS(osiResp, osvResp, serviceGrops, req.Begin, req.End, req.ForAbonent)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetPaymentsReport(ctx context.Context, req *model.PaymentsReportRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	paymentsResp, err := c.coreRepo.GetPayments(ctx, req)
	if err != nil {
		return nil, err
	}

	sort.Slice(paymentsResp, func(i, j int) bool {
		return paymentsResp[i].Dt < paymentsResp[j].Dt
	})

	rep := excel.NewPaymentsXLS(osiResp, paymentsResp, req.Begin, req.End)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetPaymentOrdersReport(ctx context.Context, req *model.PaymentOrdersRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	svodResp, err := c.coreRepo.GetPaymentOrders(ctx, req)
	if err != nil {
		return nil, err
	}

	sort.Slice(svodResp, func(i, j int) bool {
		return time.Time(svodResp[i].Date).Before(time.Time(svodResp[j].Date))
	})

	rep := excel.NewSvodXLS(osiResp, svodResp, time.Time(req.Begin), time.Time(req.End))
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetAbonentOSVReport(ctx context.Context, req *model.AbonentOSVReportRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	osvResp, err := c.coreRepo.GetAbonentOSVAll(ctx, req.AbonentID)
	if err != nil {
		return nil, err
	}

	abonentInfoResp, err := c.coreRepo.GetAbonentInfo(ctx, req.AbonentID)
	if err != nil {
		return nil, err
	}

	// fixesResp, err := c.coreRepo.GetFixes(ctx, &model.FixesRequest{
	// 	ID:    req.ID,
	// 	Begin: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
	// 	End:   time.Now(),
	// })

	// osvResp = common.MapFixesToOSV(abonentInfoResp.Flat, osvResp, fixesResp)

	rep := excel.NewSaldoAbonentXLS(osiResp, osvResp, abonentInfoResp, req.Flat)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetFixesReport(ctx context.Context, req *model.FixesRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	fixesResp, err := c.coreRepo.GetFixes(ctx, req)
	if err != nil {
		return nil, err
	}

	sort.Slice(fixesResp, func(i, j int) bool {
		return fixesResp[i].Dt < fixesResp[j].Dt
	})

	rep := excel.NewFixesXLS(osiResp, fixesResp, req.Begin, req.End)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetOSIAbonentsReport(ctx context.Context, req *model.OSIAbonentsRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	abonentsResp, err := c.coreRepo.GetOSIAbonents(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	sort.Slice(abonentsResp, func(i, j int) bool {
		return common.LessAlphaNum(abonentsResp[i].Flat, abonentsResp[j].Flat)
	})

	rep := excel.NewAbonentsXLS(osiResp, abonentsResp)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

func (c *appService) GetDebtsReport(ctx context.Context, req *model.OSVReportRequest) (report.Report, error) {
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	osvResp, err := c.coreRepo.GetOSVForDebtors(ctx, req)
	if err != nil {
		return nil, err
	}

	serviceGrops, err := c.coreRepo.GetServiceGroups(ctx)
	if err != nil {
		return nil, err
	}

	rep := excel.NewDebtsXLS(osiResp, osvResp, serviceGrops, req.Begin, req.End)
	if err := rep.Render(); err != nil {
		return nil, err
	}

	return rep, nil
}

// func (c *appService) GetAccountsReport(ctx context.Context, req *model.AccountsReportRequest) (*model.ReportResult, error) {
// 	action := "GetAccountReports"
// 	report, err := c.coreRepo.GetAccountReports(ctx, req.Path.ID)
// 	if err != nil {
// 		return nil, errors.Wrap(err, action)
// 	}

// 	osiInfo, err := c.coreRepo.GetOSIInfo(ctx, report.UnionType.ID)
// 	if err != nil {
// 		return nil, errors.Wrap(err, action)
// 	}
// 	report.Calc(req.Query.Language, osiInfo.Fio)
// 	t := template.New("GetAccountReports").Funcs(template.FuncMap{
// 		"formatAmount":  common.FormatAmount,
// 		"getDateBegin":  report.GetDateBegin,
// 		"getDateEnd":    report.GetDateEnd,
// 		"getReportDate": report.GetReportDate,
// 	})
// 	content := string(store.GetAccountsReportTemplate())
// 	t, err = t.Parse(content)
// 	if err != nil {
// 		return nil, errors.Wrap(err, action)
// 	}

// 	var buf bytes.Buffer
// 	t.Execute(&buf, &report)
// 	html := base64.StdEncoding.EncodeToString(buf.Bytes())

// 	no := false
// 	var margin uint = 5
// 	zoom := 0.7
// 	settings := model.PDFSettings{
// 		Grayscale:    &no,
// 		MarginLeft:   &margin,
// 		MarginTop:    &margin,
// 		MarginBottom: &margin,
// 		Zoom:         &zoom,
// 	}
// 	pdf, err := c.doc.HTMLToPDFBase64(html, &settings)
// 	if err != nil {
// 		return nil, errors.Wrap(err, action)
// 	}

// 	return &model.ReportResult{
// 		DOCBase64: pdf,
// 	}, nil
// }

func (c *appService) FillAccountsReport(ctx context.Context, req *model.FillReportRequest) (*model.ReportResult, error) {
	action := "FillAccountsReport"
	report := req.Body

	b, _ := json.Marshal(report)
	c.log.Debugf("FillAccountsReport request body: %s", string(b))

	report.Calc("ru", "")
	t := template.New("GetAccountReports").Funcs(template.FuncMap{
		"formatAmount":     common.FormatAmount,
		"getDateBegin":     report.GetDateBegin,
		"getDateEnd":       report.GetDateEnd,
		"getReportDate":    report.GetReportDate,
		"getClassTableRow": report.GetClassTableRow,
		"getClassCol3":     report.GetClassCol3,
	})
	content := string(store.GetAccountsReportTemplate())
	t, err := t.Parse(content)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	var buf bytes.Buffer
	t.Execute(&buf, &report)
	htmlRU := buf.Bytes()

	report.Calc("kz", "")
	tkz := template.New("GetAccountReportsKZ").Funcs(template.FuncMap{
		"formatAmount":     common.FormatAmount,
		"getDateBeginKZ":   report.GetDateBeginKZ,
		"getDateEndKZ":     report.GetDateEndKZ,
		"getReportDateKZ":  report.GetReportDateKZ,
		"getClassTableRow": report.GetClassTableRow,
		"getClassCol3":     report.GetClassCol3,
	})
	contentKZ := string(store.GetAccountsReportKZTemplate())
	tkz, err = tkz.Parse(contentKZ)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	var bufKZ bytes.Buffer
	tkz.Execute(&bufKZ, &report)
	htmlKZ := bufKZ.Bytes()

	no := false
	var margin uint = 5
	zoom := 0.7
	settings := model.PDFSettings{
		Grayscale:    &no,
		MarginLeft:   &margin,
		MarginTop:    &margin,
		MarginBottom: &margin,
		Zoom:         &zoom,
	}
	pdfRU, err := c.doc.HTMLToPDF(htmlRU, &settings)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	pdfKZ, err := c.doc.HTMLToPDF(htmlKZ, &settings)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	pdfs := [][]byte{pdfRU}
	pdfs = append(pdfs, pdfKZ)

	pdf, err := common.MergePDF(pdfs)
	if err != nil {
		return nil, err
	}

	return &model.ReportResult{
		DOCBase64: base64.StdEncoding.EncodeToString(pdf),
	}, nil
}
