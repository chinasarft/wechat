package message

type EventUnsubscribeMessageHandler func(*EventUnsubscribeRequest) *EventUnsubscribeResponse

var eventUnsubscribeMessageHandler EventUnsubscribeMessageHandler = nil

type EventUnsubscribeRequest struct {
	MessageHeader
	Event EventType
}

type EventUnsubscribeResponse struct {
	MessageHeader
	Content string
}

func SetEventUnsubscribeMessageHandler(f EventUnsubscribeMessageHandler) {

	eventUnsubscribeMessageHandler = f
}

func (this *EventUnsubscribeRequest) NewResponse(content string) *EventUnsubscribeResponse {

	return &EventUnsubscribeResponse{this.MessageHeader.getResponseHeader(), content}
}
