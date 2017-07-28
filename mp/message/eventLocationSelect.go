package message

type EventLocationSelectMessageHandler func(*EventLocationSelectRequest) *EventLocationSelectResponse

var eventLocationSelectMessageHandler EventLocationSelectMessageHandler = nil

//EventLocationSelect 文本消息
type EventLocationSelectRequest struct {
	MessageHeader
	Event            EventType
	EventKey         string
	SendLocationInfo SendLocationInfo
}

type EventLocationSelectResponse struct {
	MessageHeader
	Content string
}

func SetEventLocationSelectMessageHandler(f EventLocationSelectMessageHandler) {

	eventLocationSelectMessageHandler = f
}

//NewEventLocationSelect 初始化文本消息
func (this *EventLocationSelectRequest) NewResponse(content string) *EventLocationSelectResponse {

	return &EventLocationSelectResponse{this.MessageHeader.getResponseHeader(), content}
}
