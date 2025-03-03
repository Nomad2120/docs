package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

// FillNotficationDebetor  godoc
// @ID FillNotficationDebetor
// @Summary Заполнение шаблона Уведомления должнику.
// @Description   Заполнение шаблона Уведомления должнику.
// @Tags Уведомления
// @Accept json
// @Produce json
// @Param abonentId path int true "id Абонента"
// @Success 200 {object} model.FillDebetorResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/notification/debetor/{abonentId} [GET]
// @BasePath /
func (h *Handler) FillNotficationDebetor(c *gin.Context) {
	ctx, cancel := h.newContext(c, "FillNotficationDebetor", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.DebtorNotificationRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.FillDebtorNotification(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении контракта", err)
		return
	}

	h.ok(ctx, c, &model.FillDebetorResponse{Result: doc})
}

// FillNotficationNotary  godoc
// @ID FillNotficationNotary
// @Summary Заполнение шаблона Уведомления нотариусу.
// @Description   Заполнение шаблона Уведомления нотариусу.
// @Tags Уведомления
// @Accept json
// @Produce json
// @Param abonentId path int true "id Абонента"
// @Success 200 {object} model.FillNotaryApplicationResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/notification/notary/{abonentId} [GET]
// @BasePath /
func (h *Handler) FillNotficationNotary(c *gin.Context) {
	ctx, cancel := h.newContext(c, "FillNotficationNotary", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.NotaryApplicationRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.FillNotaryNotification(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении контракта", err)
		return
	}

	h.ok(ctx, c, &model.FillNotaryApplicationResponse{Result: doc})
}
