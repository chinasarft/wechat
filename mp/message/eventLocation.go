package message

type EventLocationMessageHandler func(*EventLocationRequest) *EventLocationResponse

var eventLocationMessageHandler EventLocationMessageHandler = nil

//EventLocation 文本消息
type EventLocationRequest struct {
	MessageHeader
	Event     EventType
	Latitude  string
	Longitude string
	Precision string
}

type EventLocationResponse struct {
	MessageHeader
	Content string
}

func SetEventLocationMessageHandler(f EventLocationMessageHandler) {

	eventLocationMessageHandler = f
}

//NewEventLocation 初始化文本消息
func (this *EventLocationRequest) NewResponse(content string) *EventLocationResponse {

	return &EventLocationResponse{this.MessageHeader.getResponseHeader(), content}
}
