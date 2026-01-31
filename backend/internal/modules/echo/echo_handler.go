package echo

import (
	"net/http"

	"github.com/D3rise/dchat/internal/errors"
	"github.com/D3rise/dchat/internal/modules/echo/dto"
	"github.com/D3rise/dchat/internal/modules/echo/services"
	"github.com/D3rise/dchat/internal/modules/echo/transformers"
	"github.com/gin-gonic/gin"
)

type EchoHandler struct {
	echoService *services.EchoService
}

func NewEchoHandler(echoService *services.EchoService) *EchoHandler {
	return &EchoHandler{
		echoService: echoService,
	}
}

func (e *EchoHandler) EchoTextHandler(c *gin.Context) {
	var req dto.EchoTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError(err.Error()))
	}

	result := e.echoService.EchoText(req.Text)

	c.JSON(http.StatusOK, transformers.EchoTextResultToResponse(result))
}
