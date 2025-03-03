package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

type ExcelFile io.Writer

// GetOSVReport  godoc
// @ID GetOSVReport
// @Summary Получение отчета сальдо за текущий месяц.
// @Description   Получение отчета сальдо за текущий месяц.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param begin query string true "дата начала"
// @Param end query string true "дата окончания"
// @Param forAbonent query bool false "для абонента"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/osv/{id} [GET]
// @BasePath /
func (h *Handler) GetOSVReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetOSVReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.OSVReportRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.ID = pathReq.ID

	resp, err := h.service.GetOSVReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetPaymentsReport  godoc
// @ID GetPaymentsReport
// @Summary Получение реестра платежей за период.
// @Description   Получение реестра платежей за период.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param begin query string true "дата начала"
// @Param end query string true "дата окончания"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/payments/{id} [GET]
// @BasePath /
func (h *Handler) GetPaymentsReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetPaymentsReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.PaymentsReportRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.ID = pathReq.ID

	resp, err := h.service.GetPaymentsReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetPaymentOrdersReport  godoc
// @ID GetPaymentOrdersReport
// @Summary Свод платежных поручений по ОСИ за период.
// @Description  Свод платежных поручений по ОСИ за период.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param begin query string true "дата начала"
// @Param end query string true "дата окончания"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/payments/orders/{id} [GET]
// @BasePath /
func (h *Handler) GetPaymentOrdersReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetPaymentOrdersReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.PaymentOrdersRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.ID = pathReq.ID

	resp, err := h.service.GetPaymentOrdersReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetAbonentOSVReport  godoc
// @ID GetAbonentOSVReport
// @Summary Получение отчета сальдо по абоненту за все периоды.
// @Description   Получение отчета сальдо по абоненту за все периоды.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param abonentId path int true "id Абонента"
// @Param flat query string true "Номер помещения"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/osv/{id}/abonent/{abonentId} [GET]
// @BasePath /
func (h *Handler) GetAbonentOSVReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetAbonentOSVReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.AbonentOSVReportRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}

	resp, err := h.service.GetAbonentOSVReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetFixesReport  godoc
// @ID GetFixesReport
// @Summary Корректировки по всем абонентам и услугам ОСИ за период.
// @Description   Корректировки по всем абонентам и услугам ОСИ за период.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param begin query string true "дата начала"
// @Param end query string true "дата окончания"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/fixes/{id} [GET]
// @BasePath /
func (h *Handler) GetFixesReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetFixesReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.FixesRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.ID = pathReq.ID

	resp, err := h.service.GetFixesReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetOSIAbonentsReport  godoc
// @ID GetOSIAbonentsReport
// @Summary Получить таблицу помещений.
// @Description   Получить таблицу помещений.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/abonents/{id} [GET]
// @BasePath /
func (h *Handler) GetOSIAbonentsReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetOSIAbonentsReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.OSIAbonentsRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	req.ID = pathReq.ID

	resp, err := h.service.GetOSIAbonentsReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetDebtsReport  godoc
// @ID GetDebtsReport
// @Summary Получение отчета по долгам за текущий месяц.
// @Description   Получение отчета  по долгам за текущий месяц.
// @Tags Отчеты
// @Accept json
// @Accept octet-stream
// @Produce  json
// @Produce  octet-stream
// @Param id path int true "id OSI"
// @Param begin query string true "дата начала"
// @Param end query string true "дата окончания"
// @Success 200 {string} model.ReportResponse "File Base64"
// @Success 200 {file} file
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/debts/{id} [GET]
// @BasePath /
func (h *Handler) GetDebtsReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetDebtsReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var (
		pathReq model.ReportRequest
		req     model.OSVReportRequest
	)
	if err := c.ShouldBindUri(&pathReq); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}
	req.ID = pathReq.ID

	resp, err := h.service.GetDebtsReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.reportOK(ctx, c, resp)
}

// GetAccountsMonthlyReport  godoc
// @ID GetAccountsMonthlyReport
// @Summary Формирование ежемесячного отчета по счетам.
// @Description Формирование ежемесячного отчета по счетам.
// @Tags Отчеты
// @Accept json
// @Produce  json
// @Param language query string false "язык отчета"
// @Param Model body model.AccountReportsResult true "модель отчета"
// @Success 200 {object} model.ReportResponse "File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/report/accounts/monthly [POST]
// @BasePath /
func (h *Handler) GetAccountsMonthlyReport(c *gin.Context) {
	ctx, cancel := h.newContext(c, "GetAccountsMonthlyReport", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.FillReportRequest

	if err := c.ShouldBindQuery(&req.Query); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	if err := c.ShouldBindJSON(&req.Body); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные", err)
		return
	}


	resp, err := h.service.FillAccountsReport(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -3, "Ошибка при формировании отчета", err)
		return
	}

	h.ok(ctx, c, resp)
}

func (h *Handler) GetTestReport(c *gin.Context) {
	time.Sleep(120 * time.Second)
	c.String(200, "ok")
}
