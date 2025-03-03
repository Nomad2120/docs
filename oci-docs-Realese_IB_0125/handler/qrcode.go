package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

// FillQRPage  godoc
// @ID FillQRPage
// @Summary Заполнение шаблона акта.
// @Description  Заполнение шаблона акта.
// @Tags Квитанции
// @Accept json
// @Produce json
// @Param id path string true "ОСИ ИД"
// @Success 200 {object} model.FillQRPageResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/invoices/qrpage/{id} [GET]
// @BasePath /
func (h *Handler) FillQRPage(c *gin.Context) {
	ctx, cancel := h.newContext(c, "FillQRPage", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.FillQRPageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.FillQRPage(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении контракта", err)
		return
	}

	h.ok(ctx, c, &model.FillQRPageResponse{Result: doc})
}
