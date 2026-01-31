package services

type EchoTextResult struct {
	Text string
}

type EchoService struct{}

func NewEchoService() *EchoService {
	return &EchoService{}
}

func (e *EchoService) EchoText(text string) EchoTextResult {
	return EchoTextResult{Text: text}
}
