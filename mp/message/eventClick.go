package message

type EventClickMessageHandler func(*EventClickRequest) *EventClickResponse

var eventClickMessageHandler EventClickMessageHandler = nil

//Click 文本消息
type EventClickRequest struct {
	MessageHeader
	Event    EventType
	EventKey string
}

type EventClickResponse struct {
	MessageHeader
	Content string
}

func SetEventClickMessageHandler(f EventClickMessageHandler) {

	eventClickMessageHandler = f
}

//NewClick 初始化文本消息
func (this *EventClickRequest) NewResponse(content string) *EventClickResponse {

	return &EventClickResponse{this.MessageHeader.getResponseHeader(), content}
}
