package dto

type EchoTextRequest struct {
	Text string `json:"text" binding:"required"`
}

type EchoTextResponse struct {
	Text string `json:"text"`
}
