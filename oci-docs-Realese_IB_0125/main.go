package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/conv"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/crypto/kalkan"
	_ "gitlab.enterprise.qazafn.kz/oci/oci-docs/docs"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/service"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/signer"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/handler"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/middleware"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/repository"
)

func main() {
	cfg := config.New(".")
	logger := common.InitLogger(cfg)
	httpClient := common.InitHTTPClient(cfg, logger)
	crypto, err := kalkan.NewKalkanCrypto()
	if err != nil {
		logger.Fatalf("Error create crypto API: %v", err)
	}
	defer crypto.Close()

	containers := []signer.CAContainer{
		{Container: "ca/GOST_03-05-2024-11_35.p12", Password: cfg.OsiStorePass, Alias: "OSI_STORE"},
		{Container: "ca/GOSTKNCA_2a8995ec50c39fe1c4d51d67c8d0bc23510065dc.p12", Password: cfg.QAZAFNStorePass, Alias: "QAZAFN"},
		{Container: "ca/GOSTKNCA_374ccee7f5e9b035c538aa54b38d6609ffe7cf4f.p12", Password: cfg.NurtauStorePass, Alias: "NURTAU"},
	}
	signer := signer.NewSigner(logger, crypto, containers)

	converter, err := conv.NewWKDocConverter()
	if err != nil {
		logger.Fatalf("Error create doc converter: %v", err)
	}
	coreRepo := repository.NewCoreHTTPRepo(cfg, logger, httpClient)
	earRepo := repository.NewEarHTTPRepo(cfg, logger, httpClient)
	appService := service.NewAppService(cfg, logger, coreRepo, converter, signer, earRepo)

	h := handler.New(cfg, logger, appService)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"origin", "x-requested-with", "accept", "content-type", "authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//l := r.Group("/", middleware.LogRequest(logger, cfg.LoggingRequest))
	api := r.Group("/api", middleware.LogRequest(logger, cfg.LoggingRequest))
	api.GET("/doc", h.GetDoc)
	api.POST("/contract", h.FillContract)
	api.POST("/contract/:id/sign", h.SignContract)
	api.POST("/contract/osi/:id/sign", h.SignOsiContract)
	api.GET("/report/osv/:id", h.GetOSVReport)
	api.GET("/report/osv/:id/abonent/:abonentId", h.GetAbonentOSVReport)
	api.GET("/report/payments/:id", h.GetPaymentsReport)
	api.GET("/report/payments/orders/:id", h.GetPaymentOrdersReport)
	api.GET("/report/fixes/:id", h.GetFixesReport)
	api.GET("/report/abonents/:id", h.GetOSIAbonentsReport)
	api.GET("/report/debts/:id", h.GetDebtsReport)
	api.POST("/report/accounts/monthly", h.GetAccountsMonthlyReport)
	//api.GET("/report/test", h.GetTestReport)

	api.GET("/act/work-completion/:id", h.FillAct)
	api.PUT("/act/work-completion/:id/sign", h.SaveAct)

	api.GET("/notification/debetor/:abonentId", h.FillNotficationDebetor)
	api.GET("/notification/notary/:abonentId", h.FillNotficationNotary)

	api.GET("/invoices/qrpage/:id", h.FillQRPage)

	pki := r.Group("/pki", middleware.LogRequest(logger, cfg.LoggingRequest))
	pki.POST("/sign/wsse", h.SignWSSE)

	r.GET("/ws", h.WS)

	docs := r.Group("/docs", middleware.LogRequest(logger, cfg.LoggingRequest))
	docs.POST("/html2pdf", h.Html2Pdf)

	// api.Use(cors.Default())
	logger.Infof("App starting on port %d", cfg.ServerPort)
	if err := r.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		logger.Errorf(err.Error())
	}
}
