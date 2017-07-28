package message

type EventViewMessageHandler func(*EventViewRequest) *EventViewResponse

var eventViewMessageHandler EventViewMessageHandler = nil

//View 文本消息
type EventViewRequest struct {
	MessageHeader
	Event    EventType
	EventKey string
	MenuId   string
}

type EventViewResponse struct {
	MessageHeader
	Content string
}

func SetEventViewMessageHandler(f EventViewMessageHandler) {

	eventViewMessageHandler = f
}

//NewView 初始化文本消息
func (this *EventViewRequest) NewResponse(content string) *EventViewResponse {

	return &EventViewResponse{this.MessageHeader.getResponseHeader(), content}
}
