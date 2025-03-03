package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	config "gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"

	"go.uber.org/zap"
)

type coreHTTPRepo struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	httpClient *http.Client
}

// NewCoreHTTPRepo -
func NewCoreHTTPRepo(cfg *config.Config, log *zap.SugaredLogger, httpClient *http.Client) Core {
	return &coreHTTPRepo{cfg: cfg, log: log, httpClient: httpClient}
}

func (c *coreHTTPRepo) SaveContract(ctx context.Context, req *model.SignContractRequest) (string, error) {
	url := fmt.Sprintf("%s/api/Registrations/%d/sign?extension=%s", c.cfg.CoreURL, req.ID, req.Extension)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp model.BaseResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "PUT", url, []byte(`"`+req.DocBase64+`"`), headers, &resp)
	if err != nil {
		return "", errors.Wrapf(err, "SaveContract. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("SaveContract. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("SaveContract. %d - %s", resp.Code, resp.Message)
	}
	return string(body), nil
}

func (c *coreHTTPRepo) GetRegistration(ctx context.Context, id int) (*model.RegistrationInfo, error) {
	url := fmt.Sprintf("%s/api/Registrations/%d", c.cfg.CoreURL, id)
	var resp model.RegistrationResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetRegistration. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetRegistration. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetRegistration. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetServiceGroups(ctx context.Context) ([]model.ServiceGroup, error) {
	var resp model.GetServiceGroupsResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", c.cfg.CoreURL+"/api/Catalogs/service-groups", nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetServiceGroups. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetServiceGroups. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetServiceGroups. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetOSVCurrentMonth(ctx context.Context, id int) (*model.OSVResult, error) {
	url := fmt.Sprintf("%s/api/Transactions/osv-current-month/%d", c.cfg.CoreURL, id)
	var resp model.OSVResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetOSVCurrentMonth. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetOSVCurrentMonth. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetOSVCurrentMonth. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetOSVByPeriod(ctx context.Context, req *model.OSVReportRequest) (*model.OSVResult, error) {
	url := fmt.Sprintf("%s/api/Transactions/osv/%d?dateBegin=%s&dateEnd=%s", c.cfg.CoreURL, req.ID, req.Begin.Format(time.RFC3339), req.End.Format(time.RFC3339))
	var resp model.OSVResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetOSVByPeriod. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetOSVByPeriod. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetOSVByPeriod. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetAbonentOSVAll(ctx context.Context, abonentId int) ([]model.AbonentOSVResult, error) {
	url := fmt.Sprintf("%s/api/Transactions/saldo-on-all-periods/%d", c.cfg.CoreURL, abonentId)
	var resp model.AbonentOSVResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetAbonentOSVAll. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetAbonentOSVAll. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetAbonentOSVAll. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetOSVForDebtors(ctx context.Context, req *model.OSVReportRequest) (*model.OSVResult, error) {
	url := fmt.Sprintf("%s/api/Transactions/osv-for-debtors/%d?dateBegin=%s&dateEnd=%s", c.cfg.CoreURL, req.ID, req.Begin.Format(time.RFC3339), req.End.Format(time.RFC3339))
	var resp model.OSVResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetOSVForDebtors. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetOSVForDebtors. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetOSVForDebtors. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetOSIInfo(ctx context.Context, id int) (*model.OSIResult, error) {
	url := fmt.Sprintf("%s/api/Osi/%d", c.cfg.CoreURL, id)
	var resp model.OSIResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetOSIInfo. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetOSIInfo. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetOSIInfo. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetAllOSI(ctx context.Context) ([]model.OSIResult, error) {
	url := c.cfg.CoreURL + "/api/Osi/all"
	var resp model.AllOSIResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetAllOSI. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetAllOSI. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetAllOSI. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil

}

func (c *coreHTTPRepo) GetPayments(ctx context.Context, req *model.PaymentsReportRequest) ([]model.PaymentInfo, error) {
	url := fmt.Sprintf("%s/api/Transactions/payments/%d?dateBegin=%s&dateEnd=%s", c.cfg.CoreURL, req.ID, req.Begin.Format(time.RFC3339), req.End.Format(time.RFC3339))
	var resp model.PaymentsResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetPayments. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetPayments. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetPayments. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetFixes(ctx context.Context, req *model.FixesRequest) ([]model.FixesInfo, error) {
	url := fmt.Sprintf("%s/api/Transactions/fixes/%d?dateBegin=%s&dateEnd=%s", c.cfg.CoreURL, req.ID, req.Begin.Format("2006-01-02T15:04:05"), req.End.Format("2006-01-02T15:04:05"))
	var resp model.FixesResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetFixes. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetFixes. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetFixes. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetAct(ctx context.Context, id int) (*model.GetActResult, error) {
	url := fmt.Sprintf("%s/api/Acts/%d", c.cfg.CoreURL, id)
	var resp model.GetActResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetAct. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetAct. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetAct. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) SignAct(ctx context.Context, req *model.SignActRequest) error {
	url := fmt.Sprintf("%s/api/Acts/%d/sign?extension=%s", c.cfg.CoreURL, req.ID, req.Extension)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp model.BaseResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "PUT", url, []byte(`"`+req.DocBase64+`"`), headers, &resp)
	if err != nil {
		return errors.Wrapf(err, "SignAct. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return fmt.Errorf("SignAct. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return fmt.Errorf("SignAct. %d - %s", resp.Code, resp.Message)
	}
	return nil
}

func (c *coreHTTPRepo) UnsignAct(ctx context.Context, id int) error {
	url := fmt.Sprintf("%s/api/Acts/%d/unsign", c.cfg.CoreURL, id)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp model.BaseResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "PUT", url, nil, headers, &resp)
	if err != nil {
		return errors.Wrapf(err, "UnsignAct. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return fmt.Errorf("UnsignAct. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return fmt.Errorf("UnsignAct. %d - %s", resp.Code, resp.Message)
	}
	return nil
}

func (c *coreHTTPRepo) GetSignedActs(ctx context.Context, id int) ([]model.SignedAct, error) {
	url := fmt.Sprintf("%s/api/Osi/%d/signed-acts", c.cfg.CoreURL, id)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp struct {
		model.BaseResponse
		Result []model.SignedAct
	}
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, headers, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetSignedActs. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetSignedActs. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetSignedActs. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetActDocs(ctx context.Context, id int) ([]model.ActDoc, error) {
	url := fmt.Sprintf("%s/api/Acts/%d/docs", c.cfg.CoreURL, id)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp struct {
		model.BaseResponse
		Result []model.ActDoc
	}
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, headers, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetActDocs. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetActDocs. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetActDocs. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetScan(ctx context.Context, id int) (string, error) {
	url := fmt.Sprintf("%s/api/Scans/%d", c.cfg.CoreURL, id)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp struct {
		model.BaseResponse
		Result string
	}
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, headers, &resp)
	if err != nil {
		return "", errors.Wrapf(err, "GetScan. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("GetScan. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return "", fmt.Errorf("GetScan. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetPaymentOrders(ctx context.Context, req *model.PaymentOrdersRequest) ([]model.PaymentOrder, error) {
	url := fmt.Sprintf("%s/api/PaymentOrders/svod/%d?dateBegin=%s&dateEnd=%s", c.cfg.CoreURL, req.ID, req.Begin.Format(time.RFC3339), req.End.Format(time.RFC3339))
	var resp model.PaymentOrdersResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetPaymentOrders. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetPaymentOrders. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetPaymentOrders. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) GetDebtorNotification(ctx context.Context, req *model.DebtorNotificationRequest) (*model.DebtorNotificationResult, error) {
	url := fmt.Sprintf("%s/api/PastDebts/debtor-notification?abonentId=%d", c.cfg.CoreURL, req.AbonentID)
	var resp model.DebtorNotificationResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetDebtorNotification. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetDebtorNotification. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetDebtorNotification. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetNotaryApplication(ctx context.Context, req *model.NotaryApplicationRequest) (*model.NotaryApplicationResult, error) {
	url := fmt.Sprintf("%s/api/PastDebts/notary-application?abonentId=%d", c.cfg.CoreURL, req.AbonentID)
	var resp model.NotaryApplicationResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetNotaryApplication. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetNotaryApplication. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetNotaryApplication. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetOSIAbonents(ctx context.Context, osiID int) ([]model.OSIAbonent, error) {
	url := fmt.Sprintf("%s/api/Osi/%d/abonents", c.cfg.CoreURL, osiID)
	var resp model.OSIAbonentsResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetOSIAbonents. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetOSIAbonents. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetOSIAbonents. %d - %s", resp.Code, resp.Message)
	}
	return resp.Result, nil
}

func (c *coreHTTPRepo) AddDoc(ctx context.Context, id int, req *model.AddDocRequest) (*model.AddDocResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "AddDoc. invalid request")
	}
	url := fmt.Sprintf("%s/api/Osi/%d/docs", c.cfg.CoreURL, id)
	headers := map[string]string{"Content-Type": "application/json"}
	var resp model.AddDocResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "POST", url, reqBody, headers, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "AddDoc. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("AddDoc. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("AddDoc. %d - %s", resp.Code, resp.Message)
	}
	return &resp, nil
}

func (c *coreHTTPRepo) GetAbonentInfo(ctx context.Context, abonentID int) (*model.OSIAbonent, error) {
	url := fmt.Sprintf("%s/api/Abonents/%d", c.cfg.CoreURL, abonentID)
	var resp model.AbonentsInfoResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetAbonentInfo. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetAbonentInfo. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetAbonentInfo. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}

func (c *coreHTTPRepo) GetAccountReports(ctx context.Context, id int) (*model.AccountReportsResult, error) {
	url := fmt.Sprintf("%s/api/AccountReports/%d/monthly-form-data", c.cfg.CoreURL, id)
	var resp model.AccountReportsResponse
	status, body, err := common.RequestJSON(ctx, c.httpClient, "GET", url, nil, nil, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "GetAccountReports. Request error. Response Body: %s", string(body))
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GetAccountReports. Invalid Response Status Code: %d", status)
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("GetAccountReports. %d - %s", resp.Code, resp.Message)
	}
	return &resp.Result, nil
}
