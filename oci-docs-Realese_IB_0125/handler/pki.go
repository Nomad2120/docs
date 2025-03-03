package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
)

const DefaultStorage string = "QAZAFN"

// SignWSSE  godoc
// @ID SignWSSE
// @Summary Подписание WSSE.
// @Description  Подписание WSSE.
// @Tags PKI
// @Accept text/plain
// @Produce text/plain
// @Param signNodeId query string true "ID элемента"
// @Param alias query string true "Имя хранилища сертификатов"
// @Param data body string true "XML plain text"
// @Success 200 {string} data "подписанный xml"
// @Failure 422 {object} model.ErrorResponse
// @Router /pki/sign/wsse [POST]
// @BasePath /
func (h *Handler) SignWSSE(c *gin.Context) {
	ctx, cancel := h.newContext(c, "SignWSSE", h.cfg.CoreTimeout.Duration)
	defer cancel()

	var req model.SignWSSERequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -1, "Неверные входные данные", err)
		return
	}

	if req.Alias == "" {
		req.Alias = DefaultStorage //"KAZAFN_STORE"
	}
	req.Alias = strings.ToUpper(req.Alias)

	reqBody, err := c.GetRawData()
	if err != nil {
		h.fail(ctx, c, http.StatusBadRequest, -2, "Не указано тело запроса", err)
		return
	}
	req.Data = string(reqBody)

	sign, err := h.service.SignWSSE(ctx, &req)
	if err != nil {
		h.fail(ctx, c, 422, -3, err.Error(), nil)
		return
	}

	c.String(http.StatusOK, sign)
}
