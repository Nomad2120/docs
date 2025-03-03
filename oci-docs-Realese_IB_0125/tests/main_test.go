package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	httpmock "gitlab.enterprise.qazafn.kz/oci/oci-docs/common/mocks"
	config "gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	conv "gitlab.enterprise.qazafn.kz/oci/oci-docs/conv"
	handler "gitlab.enterprise.qazafn.kz/oci/oci-docs/handler"
	repository "gitlab.enterprise.qazafn.kz/oci/oci-docs/repository"
	service "gitlab.enterprise.qazafn.kz/oci/oci-docs/service"
)

// func TestGetAccountsMonthlyReport(t *testing.T) {
// 	cfg := config.New("..")
// 	logger := common.InitLogger(cfg)

// 	osiResp := common.LoadFixture(t, "../testdata/osi_info_response.json")
// 	reportResp := common.LoadFixture(t, "../testdata/accounts_report_response.json")

// 	header := make(http.Header)
// 	header.Set("content-type", "application/json")
// 	responses := []httpmock.ResponseHandler{
// 		{
// 			Response: &http.Response{
// 				StatusCode: 200,
// 				Body:       io.NopCloser(bytes.NewBuffer(reportResp)),
// 				Header:     header,
// 			},
// 			Handler: nil,
// 		},
// 		{
// 			Response: &http.Response{
// 				StatusCode: 200,
// 				Body:       io.NopCloser(bytes.NewBuffer(osiResp)),
// 				Header:     header,
// 			},
// 			Handler: nil,
// 		},
// 	}

// 	rt := httpmock.NewRoundTripper(responses)
// 	httpClient := &http.Client{Transport: rt}

// 	converter, err := conv.NewWKDocConverter()
// 	if err != nil {
// 		logger.Fatalf("Error create doc converter: %v", err)
// 	}
// 	coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
// 	earRepo := repository.NewEarHTTPRepo(cfg, logger, httpClient)
// 	appService := service.NewAppService(cfg, logger, coreRepo, converter, nil, earRepo)

// 	h := handler.New(cfg, logger, appService)

// 	gin.SetMode(gin.ReleaseMode)
// 	gin.DefaultWriter = ioutil.Discard

// 	r := gin.Default()

// 	api := r.Group("/api")
// 	api.GET("/report/accounts/:id/monthly", h.GetAccountsMonthlyReport)

// 	w := httpmock.DoRequest(r, "GET", "/api/report/accounts/1/monthly?language=ru", "", "", nil, nil)
// 	assert.Equal(t, http.StatusOK, w.Code, "HTTP Status should be 200")

// 	body := w.Body.Bytes()
// 	println(string(body))
// }

func TestGetAccountsMonthlyReport(t *testing.T) {
	cfg := config.New("..")
	logger := common.InitLogger(cfg)

	converter, err := conv.NewWKDocConverter()
	if err != nil {
		logger.Fatalf("Error create doc converter: %v", err)
	}
	coreRepo := repository.NewCoreHTTPRepo(cfg, logger, nil)
	earRepo := repository.NewEarHTTPRepo(cfg, logger, nil)
	appService := service.NewAppService(cfg, logger, coreRepo, converter, nil, earRepo)

	h := handler.New(cfg, logger, appService)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	api := r.Group("/api")
	api.GET("/report/accounts/:id/monthly", h.GetAccountsMonthlyReport)

	reportResp := common.LoadFixture(t, "../testdata/accounts_report_response.json")
	w := httpmock.DoRequest(r, "GET", "/api/report/accounts/1/monthly?language=ru", "", "", nil, io.NopCloser(bytes.NewBuffer(reportResp)))
	assert.Equal(t, http.StatusOK, w.Code, "HTTP Status should be 200")

	body := w.Body.Bytes()
	println(string(body))
}
