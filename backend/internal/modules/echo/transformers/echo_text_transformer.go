package transformers

import (
	"github.com/D3rise/dchat/internal/modules/echo/dto"
	"github.com/D3rise/dchat/internal/modules/echo/services"
)

func EchoTextResultToResponse(result services.EchoTextResult) dto.EchoTextResponse {
	return dto.EchoTextResponse{Text: result.Text}
}
