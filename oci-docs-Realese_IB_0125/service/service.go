package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common/docx"
	config "gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	conv "gitlab.enterprise.qazafn.kz/oci/oci-docs/conv"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	report "gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
	repository "gitlab.enterprise.qazafn.kz/oci/oci-docs/repository"
	signer "gitlab.enterprise.qazafn.kz/oci/oci-docs/signer"
	store "gitlab.enterprise.qazafn.kz/oci/oci-docs/template"
	"go.uber.org/zap"
)

type AppService interface {
	GetDoc(context.Context, *model.DocRequest) (string, error)
	Html2Pdf(context.Context, *model.Html2PdfRequest) (string, error)
	FillContract(context.Context, *model.FillContractRequest) (*model.FillContractResult, error)
	SignContract(context.Context, *model.SignContractRequest) (*model.SignContractResponse, error)
	FillAct(context.Context, *model.FillActRequest) (*model.FillActResult, error)
	FillDebtorNotification(context.Context, *model.DebtorNotificationRequest) (*model.FillDebetorResult, error)
	FillNotaryNotification(context.Context, *model.NotaryApplicationRequest) (*model.FillNotaryApplicationResult, error)

	GetOSVReport(context.Context, *model.OSVReportRequest) (report.Report, error)
	GetAbonentOSVReport(ctx context.Context, req *model.AbonentOSVReportRequest) (report.Report, error)
	GetPaymentsReport(context.Context, *model.PaymentsReportRequest) (report.Report, error)
	GetPaymentOrdersReport(context.Context, *model.PaymentOrdersRequest) (report.Report, error)
	GetFixesReport(context.Context, *model.FixesRequest) (report.Report, error)
	GetOSIAbonentsReport(context.Context, *model.OSIAbonentsRequest) (report.Report, error)
	GetDebtsReport(context.Context, *model.OSVReportRequest) (report.Report, error)
	//GetAccountsReport(context.Context, *model.AccountsReportRequest) (*model.ReportResult, error)
	FillAccountsReport(context.Context, *model.FillReportRequest) (*model.ReportResult, error)

	SaveAct(ctx context.Context, req *model.SignActRequest) error
	SignWSSE(ctx context.Context, req *model.SignWSSERequest) (string, error)

	FillQRPage(context.Context, *model.FillQRPageRequest) (*model.FillQRPageResult, error)
	SignOsiContract(context.Context, *model.SignOsiContractRequest) (*model.AddDocResponse, error)
}

type appService struct {
	cfg      *config.Config
	log      *zap.SugaredLogger
	doc      conv.DocConverter
	signer   signer.Signer
	coreRepo repository.Core
	earRepo  repository.Ear
}

func NewAppService(cfg *config.Config, log *zap.SugaredLogger, coreRepo repository.Core, doc conv.DocConverter, signer signer.Signer,
	earRepo repository.Ear) AppService {
	return &appService{cfg: cfg, log: log, doc: doc, coreRepo: coreRepo, signer: signer, earRepo: earRepo}
}

// GetDoc -
func (c *appService) GetDoc(ctx context.Context, req *model.DocRequest) (string, error) {
	enc := base64.StdEncoding.EncodeToString(store.GetContractTemplate())

	pdf, err := c.doc.HTMLToPDFBase64(enc, nil)
	if err != nil {
		return "", errors.Wrap(err, "html2pdf")
	}

	signDoc, err := c.signer.SignCMSBase64("OSI_STORE", pdf)
	if err != nil {
		return "", errors.Wrap(err, "Ошибка подписи")
	}

	return signDoc, nil
}

// Html2Pdf -
func (c *appService) Html2Pdf(ctx context.Context, req *model.Html2PdfRequest) (string, error) {
	pdf, err := c.doc.HTMLToPDFBase64(req.HTMLBase64, &req.Settings)
	if err != nil {
		return "", errors.Wrap(err, "html2pdf")
	}
	return pdf, nil
}

func (c *appService) FillContract(ctx context.Context, req *model.FillContractRequest) (*model.FillContractResult, error) {
	regInfo, err := c.coreRepo.GetRegistration(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "FillContract")
	}

	if regInfo.AddressKZ == "" {
		buildingResp, err := c.earRepo.GetBuildingInfo(ctx, regInfo.AtsId, regInfo.AddressRegistryId)
		if err != nil {
			regInfo.AddressKZ = regInfo.Address
		} else {
			regInfo.AddressKZ = buildingResp.ShortPathKaz
		}
	}
	req.UnionTypeRu = regInfo.UnionTypeRu
	req.UnionTypeKz = regInfo.UnionTypeKz
	req.Site = "https://eosi.kz/"
	t := template.New("FillContract").Funcs(template.FuncMap{
		"formatAmount": common.FormatAmount,
		"strAmount":    common.Num2Str,
		"strAmountKaz": common.Num2StrKaz,
		"signDate":     req.SignDate,
	})
	tmpl, err := t.Parse(string(store.GetContractTemplate()))
	if err != nil {
		return nil, errors.Wrap(err, "FillContract")
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, req)
	html := base64.StdEncoding.EncodeToString(buf.Bytes())

	pdf, err := c.doc.HTMLToPDFBase64(html, nil)
	if err != nil {
		return nil, errors.Wrap(err, "FillContract")
	}

	return &model.FillContractResult{
		HTMLBase64: html,
		PDFBase64:  pdf,
	}, nil
}

func (c *appService) SignContract(ctx context.Context, req *model.SignContractRequest) (*model.SignContractResponse, error) {
	iin, cert, err := c.signer.GetIIN(req.DocBase64, 1)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при получении сертификата")
	}
	regResp, err := c.coreRepo.GetRegistration(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при получении данных регистрации")
	}

	if strings.EqualFold(regResp.StateCode, "SIGNED") {
		return nil, errors.New("Заявка уже была подписана ранее.")
	}

	if cert.NotAfter.Before(time.Now()) {
		return nil, errors.New("Срок действия сертификат истек " + cert.NotAfter.Format("02.01.2006"))
	}

	if iin != regResp.IDN {
		return nil, fmt.Errorf("ИИН/БИН %s регистрационных данных не соответствует IIN%s сертификата", regResp.IDN, iin)
	}

	sigDoc, err := c.signer.SignCMSBase64("OSI_STORE", req.DocBase64)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при подписании контракта")
	}

	req.DocBase64 = sigDoc
	req.Extension = "pdf"

	doc, err := c.coreRepo.SaveContract(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при сохранении контракта")
	}

	return &model.SignContractResponse{Result: doc}, nil
}

func (c *appService) FillAct(ctx context.Context, req *model.FillActRequest) (*model.FillActResult, error) {
	act, err := c.coreRepo.GetAct(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "FillAct")
	}
	t := template.New("FillAct").Funcs(template.FuncMap{
		"dateContract":  act.GetDateContract,
		"discount":      act.GetDiscount,
		"total":         act.GetTotal,
		"oldDiscount":   act.GetOldDiscount,
		"oldTotal":      act.GetOldTotal,
		"totalQuantity": act.GetTotalQuantity,
		"formatAmount":  common.FormatAmount,
		"fullName":      act.GetFullName,
		"getPeriod":     act.GetPeriod,
	})
	content := ""
	if len(act.ActItems) > 0 {
		content = string(store.GetActTemplate())
	} else {
		content = string(store.GetActOldTemplate())
	}
	t, err = t.Parse(content)
	if err != nil {
		return nil, errors.Wrap(err, "FillAct")
	}

	var buf bytes.Buffer
	t.Execute(&buf, &act)
	html := base64.StdEncoding.EncodeToString(buf.Bytes())

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
	pdf, err := c.doc.HTMLToPDFBase64(html, &settings)
	if err != nil {
		return nil, errors.Wrap(err, "FillAct")
	}

	return &model.FillActResult{
		//	HTMLBase64: html,
		PDFBase64: pdf,
	}, nil
}

// func (c *appService) FillDebtorNotification(ctx context.Context, req *model.DebtorNotificationRequest) (*model.FillDebetorResult, error) {
// 	resp, err := c.coreRepo.GetDebtorNotification(ctx, req)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FillDebtorNotification")
// 	}

// 	t := template.New("FillDebtorNotification").Funcs(template.FuncMap{
// 		"formatAmount": common.FormatAmount,
// 		"num2Str": common.Num2Str,
// 		"getDay": resp.GetDay,
// 		"getMonth": resp.GetMonth,
// 		"getYear": resp.GetYear,
// 	})
// 	t, err = t.Parse(string(store.GetNotificationToDebitorTemplate()))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FillDebtorNotification")
// 	}

// 	var buf bytes.Buffer
// 	t.Execute(&buf, resp)
// 	html := base64.StdEncoding.EncodeToString(buf.Bytes())

// 	no := false
// 	var margin uint = 5
// 	zoom := 1.0
// 	settings := model.PDFSettings{
// 		Grayscale:    &no,
// 		MarginLeft:   &margin,
// 		MarginTop:    &margin,
// 		MarginBottom: &margin,
// 		Zoom:         &zoom,
// 	}
// 	pdf, err := c.doc.HTMLToPDFBase64(html, &settings)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "FillDebtorNotification")
// 	}

// 	return &model.FillDebetorResult{
// 		PDFBase64: pdf,
// 	}, nil

// }

func (c *appService) FillDebtorNotification(ctx context.Context, req *model.DebtorNotificationRequest) (*model.FillDebetorResult, error) {
	action := "FillDebtorNotification"
	resp, err := c.coreRepo.GetDebtorNotification(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	osiInfo, err := c.coreRepo.GetOSIInfo(ctx, resp.OsiID)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	if !strings.HasPrefix(resp.OsiName, osiInfo.UnionTypeRu) {
		resp.UnionTypeRu = osiInfo.UnionTypeRu
		resp.UnionTypeKz = osiInfo.UnionTypeKz
	}
	//resp.OsiName = strings.TrimPrefix(resp.OsiName, osiInfo.UnionTypeRu)

	if resp.AddressKZ == "" {
		if osiInfo.AddressKZ != "" {
			resp.AddressKZ = osiInfo.AddressKZ
		} else {
			resp.AddressKZ, err = c.getOsiAddressKz(ctx, osiInfo)
			if err != nil {
				return nil, errors.Wrap(err, action)
			}
			resp.AddressKZ += ", " + resp.Flat + " пəтер"
		}
	}

	funcMap := template.FuncMap{
		"formatAmount":   common.FormatAmount,
		"num2Str":        common.Num2Str,
		"num2StrKaz":     common.Num2StrKaz,
		"getDay":         resp.GetDay,
		"getFullDate":    resp.GetFullDate,
		"getFullDateKaz": resp.GetFullDateKZ,
		"getMonth":       resp.GetMonth,
		"getYear":        resp.GetYear,
		"getMonthKaz":    resp.GetMonthKZ,
	}
	fname := fmt.Sprintf("FillDebtorNotification_%s.docx", uuid.New().String())
	ioutil.WriteFile(fname, store.GetNotificationToDebitorTemplate(), os.ModePerm)
	defer func() {
		if err := os.Remove(fname); err != nil {
			c.log.Errorf("%s. Remove file error: %s", action, fname)
		}
	}()
	b, err := docx.RenderDocxTemplate(fname, funcMap, resp)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	data := base64.StdEncoding.EncodeToString(b)

	return &model.FillDebetorResult{
		DOCBase64: data,
	}, nil

}

func (c *appService) FillNotaryNotification(ctx context.Context, req *model.NotaryApplicationRequest) (*model.FillNotaryApplicationResult, error) {
	action := "FillNotaryNotification"
	resp, err := c.coreRepo.GetNotaryApplication(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}

	osiInfo, err := c.coreRepo.GetOSIInfo(ctx, resp.OsiID)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}
	if !strings.HasPrefix(resp.OsiName, osiInfo.UnionTypeRu) {
		resp.UnionTypeRu = osiInfo.UnionTypeRu
		resp.UnionTypeKz = osiInfo.UnionTypeKz
	}
	//resp.OsiName = strings.TrimPrefix(resp.OsiName, osiInfo.UnionTypeRu)

	if resp.OsiAddressKZ == "" {
		resp.OsiAddressKZ, err = c.getOsiAddressKz(ctx, osiInfo)
		if err != nil {
			return nil, errors.Wrap(err, action)
		}
		resp.AbonentAddressKZ = resp.OsiAddressKZ + ", " + resp.AbonentFlat + " пəтер"
	}

	if resp.AbonentIdn == "" {
		resp.AbonentIdn = "____________"
	}

	funcMap := template.FuncMap{
		"formatAmount":  common.FormatAmount,
		"formatDateNom": resp.FormatDateNom,
		"formatDate":    resp.FormatDate,
		"toLower":       strings.ToLower,
		"getDateBegin":  resp.GetDateBegin,
		"getDateEnd":    resp.GetDateEnd,
		"getDebtDate":   resp.GetDebtDate,
		"getChairman":   resp.GetChairman,
		"num2StrKaz":    common.Num2StrKaz,
	}
	fname := fmt.Sprintf("FillNotaryNotification_%s.docx", uuid.New().String())
	os.WriteFile(fname, store.GetNotaryTemplate(), os.ModePerm)
	defer func() {
		if err := os.Remove(fname); err != nil {
			c.log.Errorf("%s. Remove file error: %s", action, fname)
		}
	}()
	b, err := docx.RenderDocxTemplate(fname, funcMap, resp)
	if err != nil {
		return nil, errors.Wrap(err, action)
	}
	data := base64.StdEncoding.EncodeToString(b)

	return &model.FillNotaryApplicationResult{
		DOCBase64: data,
	}, nil

}

func (c *appService) SaveAct(ctx context.Context, req *model.SignActRequest) error {
	c.log.Debugf("SaveAct  id: %d, docbase64: %s", req.ID, req.DocBase64)
	iin, cert, err := c.signer.GetIIN(req.DocBase64, 1)
	if err != nil {
		return errors.Wrap(err, "Ошибка при получении сертификата")
	}
	act, err := c.coreRepo.GetAct(ctx, req.ID)
	if err != nil {
		return errors.Wrap(err, "Ошибка при получении данных aкта")
	}

	if cert.NotAfter.Before(time.Now()) {
		return errors.New("Срок действия сертификат истек " + cert.NotAfter.Format("02.01.2006"))
	}

	if iin != act.OsiIDN {
		return fmt.Errorf("ИИН/БИН %s регистрационных данных не соответствует IIN%s сертификата", act.OsiIDN, iin)
	}

	sigDoc, err := c.signer.SignCMSBase64("OSI_STORE", req.DocBase64)
	if err != nil {
		return errors.Wrap(err, "Ошибка при подписании акта")
	}

	req.DocBase64 = sigDoc

	if err := c.coreRepo.SignAct(ctx, req); err != nil {
		return errors.Wrap(err, "Ошибка при сохранении акта")
	}

	return nil
}

func (c *appService) SignWSSE(ctx context.Context, req *model.SignWSSERequest) (string, error) {
	sign, err := c.signer.SignWSSE(req.Alias, req.Data, req.SignNodeID)
	if err != nil {
		return "", errors.Wrap(err, "Ошибка при подписании")
	}
	n := 0
	if strings.Contains(sign, "<?xml") {
		n = strings.Index(sign, ">")
		if n == -1 {
			n = 0
		}
		if n != 0 && len(sign) > n+2 {
			if sign[n+1] == '\n' || sign[n+1] == '\r' {
				n++
			}
			if sign[n+2] == '\n' || sign[n+2] == '\r' {
				n++
			}
		}
	}
	return sign[n+1:], nil
}

func (c *appService) getOsiAddressKz(ctx context.Context, osiInfo *model.OSIResult) (string, error) {
	regInfo, err := c.coreRepo.GetRegistration(ctx, osiInfo.RegistrationID)
	if err != nil {
		return "", err
	}
	if regInfo.AddressKZ != "" {
		return regInfo.AddressKZ, nil
	}
	if regInfo.AtsId == 0 && regInfo.AddressRegistryId == 0 {
		return osiInfo.Address, nil
	}

	buildingResp, err := c.earRepo.GetBuildingInfo(ctx, regInfo.AtsId, regInfo.AddressRegistryId)
	if err != nil {
		return "", err
	}
	return buildingResp.ShortPathKaz, nil

}

func (c *appService) SignOsiContract(ctx context.Context, req *model.SignOsiContractRequest) (*model.AddDocResponse, error) {
	iin, cert, err := c.signer.GetIIN(req.Body.DocBase64, 1)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при получении сертификата")
	}
	osiResp, err := c.coreRepo.GetOSIInfo(ctx, req.Path.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при получении данных регистрации")
	}

	if cert.NotAfter.Before(time.Now()) {
		return nil, errors.New("Срок действия сертификат истек " + cert.NotAfter.Format("02.01.2006"))
	}

	if iin != osiResp.Idn {
		return nil, fmt.Errorf("ИИН/БИН %s регистрационных данных не соответствует IIN%s сертификата", osiResp.Idn, iin)
	}

	sigDoc, err := c.signer.SignCMSBase64("OSI_STORE", req.Body.DocBase64)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при подписании контракта")
	}

	resp, err := c.coreRepo.AddDoc(ctx, req.Path.ID, &model.AddDocRequest{
		DocTypeCode: "SIGNED_CONTRACT",
		Data:        sigDoc,
		Extension:   "pdf",
	})
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при сохранении контракта")
	}

	return resp, nil
}
