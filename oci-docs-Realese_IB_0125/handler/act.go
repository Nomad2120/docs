package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

// FillAct  godoc
// @ID FillAct
// @Summary Заполнение шаблона акта.
// @Description  Заполнение шаблона акта.
// @Tags Акты
// @Accept json
// @Produce json
// @Param id path string true "ид акта"
// @Success 200 {object} model.FillActResponse	"File Base64"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/act/work-completion/{id} [GET]
// @BasePath /
func (h *Handler) FillAct(c *gin.Context) {
	ctx, cancel := h.newContext(c, "FillAct", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.FillActRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	doc, err := h.service.FillAct(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 200, -2, "Ошибка при получении контракта", err)
		return
	}

	h.ok(ctx, c, &model.FillActResponse{Result: doc})
}

// SaveAct  godoc
// @ID SaveAct
// @Summary Сохранение подписанного акта.
// @Description  Сохранение подписанного акта.
// @Tags Акты
// @Accept json
// @Produce json
// @Param id path string true "ид акта"
// @Param extension query string true "расширение файла"
// @Param docBase64 body string true "документ base64"
// @Success 200 {object} model.BaseResponse "модель"
// @Failure 422 {object} model.ErrorResponse
// @Router /api/act/work-completion/{id}/sign [PUT]
// @BasePath /
func (h *Handler) SaveAct(c *gin.Context) {
	ctx, cancel := h.newContext(c, "SaveAct", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.SignActRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}
	req.Extension = c.Query("extension")
	if req.Extension == "" {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Неверные входные данные. extension required", nil)
		return
	}
	reqBody, err := c.GetRawData()
	if err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -3, "Неверные входные данные", err)
		return
	}
	req.DocBase64 = string(reqBody)
	if err := h.service.SaveAct(ctx, &req); err != nil {
		h.fail(ctx, c, 200, -4, err.Error(), nil)
		return
	}

	h.ok(ctx, c, &model.BaseResponse{})
}
