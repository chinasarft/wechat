package message

type EventSubscribeMessageHandler func(*EventSubscribeRequest) *EventSubscribeResponse

var eventSubscribeMessageHandler EventSubscribeMessageHandler = nil

type EventSubscribeRequest struct {
	MessageHeader
	Event EventType
}

type EventSubscribeResponse struct {
	MessageHeader
	Content string
}

func SetEventSubscribeMessageHandler(f EventSubscribeMessageHandler) {

	eventSubscribeMessageHandler = f
}

func (this *EventSubscribeRequest) NewResponse(content string) *EventSubscribeResponse {

	return &EventSubscribeResponse{this.MessageHeader.getResponseHeader(), content}
}
