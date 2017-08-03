package message

type EventScanMessageHandler func(*EventScanRequest) *EventScanResponse

var eventScanMessageHandler EventScanMessageHandler = nil

type EventScanRequest struct {
	MessageHeader
	Event    EventType
	EventKey string
	Ticket   string
}

type EventScanResponse struct {
	MessageHeader
	Content string
}

func SetEventScanMessageHandler(f EventScanMessageHandler) {

	eventScanMessageHandler = f
}

func (this *EventScanRequest) NewResponse(content string) *EventScanResponse {

	return &EventScanResponse{this.MessageHeader.getResponseHeader(), content}
}
