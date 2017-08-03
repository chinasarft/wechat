package message

type EventScancodeWaitmsgMessageHandler func(*EventScancodeWaitmsgRequest) *EventScancodeWaitmsgResponse

var eventScancodeWaitmsgMessageHandler EventScancodeWaitmsgMessageHandler = nil

//EventScancodeWaitmsg 文本消息
type EventScancodeWaitmsgRequest struct {
	MessageHeader
	Event        EventType
	EventKey     string
	ScanCodeInfo ScanCodeInfo
}

type EventScancodeWaitmsgResponse struct {
	MessageHeader
	Content string
}

func SetEventScancodeWaitmsgMessageHandler(f EventScancodeWaitmsgMessageHandler) {

	eventScancodeWaitmsgMessageHandler = f
}

//NewEventScancodeWaitmsg 初始化文本消息
func (this *EventScancodeWaitmsgRequest) NewResponse(content string) *EventScancodeWaitmsgResponse {

	return &EventScancodeWaitmsgResponse{this.MessageHeader.getResponseHeader(), content}
}
