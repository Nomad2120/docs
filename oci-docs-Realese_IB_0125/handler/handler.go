package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model/context/action"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model/context/auth"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model/context/processid"
	report "gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
	service "gitlab.enterprise.qazafn.kz/oci/oci-docs/service"
	"go.uber.org/zap"
)

// Handler -
type Handler struct {
	cfg      *config.Config
	log      *zap.SugaredLogger
	service  service.AppService
	conninfo string
}

// New -
func New(cfg *config.Config, log *zap.SugaredLogger, service service.AppService) *Handler {
	conninfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.CoreDBHost, cfg.CoreDBPort, cfg.CoreDBUser, cfg.CoreDBPassword, cfg.CoreDBName)

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatalf("Error connect to DB (host=%s port=%d user=%s dbname=%s): %v", cfg.CoreDBHost, cfg.CoreDBPort, cfg.CoreDBUser, cfg.CoreDBName, err)
	}
	//	db.Close()

	return &Handler{
		cfg:      cfg,
		log:      log,
		service:  service,
		conninfo: conninfo,
	}
}

func (h *Handler) newContext(c *gin.Context, act string, timeout time.Duration) (context.Context, context.CancelFunc) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx = action.NewContext(processid.NewContext(ctx), act)

	authorization, ok := c.Request.Header["Authorization"]
	if ok && len(authorization) > 0 {
		ctx = auth.NewContext(ctx, authorization[0])
	}
	return ctx, cancel
}

func (h *Handler) fail(ctx context.Context, c *gin.Context, status, code int, message string, err error) {
	fullmsg := message
	serr := ""
	if err != nil {
		serr = err.Error()
		if fullmsg != "" {
			fullmsg += ": "
		}
		fullmsg += serr
	}
	h.log.Errorw("Error handling: "+fullmsg, "status", status, "action", action.FromContext(ctx), "code", strconv.Itoa(code), "processId", processid.FromContext(ctx))
	c.JSON(status, &model.ErrorResponse{
		BaseResponse: model.BaseResponse{Code: code, Message: message},
		Result:       serr,
	})
}

func (h *Handler) ok(ctx context.Context, c *gin.Context, body interface{}) {
	h.log.Debugw("Success", "action", action.FromContext(ctx), "processId", processid.FromContext(ctx))
	c.JSON(http.StatusOK, body)
}

func (h *Handler) reportOK(ctx context.Context, c *gin.Context, resp report.Report) {
	if c.Request.Header.Get("Accept") == "application/octet-stream" {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+"OSVReport.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")

		if _, err := resp.WriteTo(c.Writer); err != nil {
			h.fail(ctx, c, 422, -4, "Ошибка при отправке отчета", err)
		}

	} else {
		var buf bytes.Buffer
		resp.WriteTo(&buf)
		h.ok(ctx, c, &model.ReportResponse{Result: base64.StdEncoding.EncodeToString(buf.Bytes())})
	}
}

// GetDoc  godoc
// @ID GetDoc
// @Summary Получение списка всех видов документов.
// @Description Получение списка всех видов документов.
// @Tags GetDoc
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id query string true "name search by id"
// @Success 200 {object} model.DocResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/doc [GET]
// @BasePath /
func (h *Handler) GetDoc(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetDoc", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.DocRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.GetDoc(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении документа: ", err)
		return
	}

	h.ok(ctx, c, &model.DocResponse{Result: doc})
}

// FillContract  godoc
// @ID FillContract
// @Summary Заполнение шаблона договора офферты.
// @Description  Заполнение шаблона договора офферты.
// @Tags Contract
// @Accept json
// @Produce json
// @Param Model body model.FillContractRequest true "модель"
// @Success 200 {object} model.FillContractResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/contract [POST]
// @BasePath /
func (h *Handler) FillContract(c *gin.Context) {
	ctx, cancel := h.newContext(c, "FillContract", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.FillContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.FillContract(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении контракта", err)
		return
	}

	h.ok(ctx, c, &model.FillContractResponse{Result: doc})
}

// SignContract  godoc
// @ID SignContract
// @Summary Подписание договора офферты.
// @Description  Подписание договора офферты.
// @Tags Contract
// @Accept json
// @Produce json
// @Param id path int true "id заявки"
// @Param docBase64 body string true "документ base64"
// @Success 200 {object} model.SignContractResponse "модель"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/contract/{id}/sign [PUT]
// @BasePath /
func (h *Handler) SignContract(c *gin.Context) {
	ctx, cancel := h.newContext(c, "SignContract", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.SignContractRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	reqBody, err := c.GetRawData()
	if err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.DocBase64 = string(reqBody)
	resp, err := h.service.SignContract(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, err.Error(), err)
		return
	}

	h.ok(ctx, c, resp)
}

// SignOsiContract  godoc
// @ID SignOsiContract
// @Summary Подписание нового договора офферты.
// @Description  Подписание нового договора офферты.
// @Tags Contract
// @Accept json
// @Produce json
// @Param id path int true "id OSI"
// @Param docBase64 body string true "документ base64"
// @Success 200 {object} model.AddDocResponse "модель"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/contract/osi/{id}/sign [POST]
// @BasePath /
func (h *Handler) SignOsiContract(c *gin.Context) {
	ctx, cancel := h.newContext(c, "SignOsiContract", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.SignOsiContractRequest
	if err := c.ShouldBindUri(&req.Path); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindJSON(&req.Body); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	resp, err := h.service.SignOsiContract(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, err.Error(), err)
		return
	}

	h.ok(ctx, c, resp)
}

// Html2Pdf  godoc
// @ID Html2Pdf
// @Summary Преобразовать html в pdf
// @Description Преобразовать html в pdf
// @Tags Docs
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Model body model.Html2PdfRequest true "модель"
// @Success 200 {object} model.DocResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /docs/html2pdf [POST]
// @BasePath /
func (h *Handler) Html2Pdf(c *gin.Context) {
	ctx, cancel := h.newContext(c, "Html2Pdf", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.Html2PdfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.Html2Pdf(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка конвертации: ", err)
		return
	}

	h.ok(ctx, c, &model.DocResponse{Result: doc})
}
