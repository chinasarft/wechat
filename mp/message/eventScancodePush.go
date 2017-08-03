package message

type EventScancodePushMessageHandler func(*EventScancodePushRequest) *EventScancodePushResponse

var eventScancodePushMessageHandler EventScancodePushMessageHandler = nil

//EventScancodePush 文本消息
type EventScancodePushRequest struct {
	MessageHeader
	Event        EventType
	EventKey     string
	ScanCodeInfo ScanCodeInfo
}

type EventScancodePushResponse struct {
	MessageHeader
	Content string
}

func SetEventScancodePushMessageHandler(f EventScancodePushMessageHandler) {

	eventScancodePushMessageHandler = f
}

//NewEventScancodePush 初始化文本消息
func (this *EventScancodePushRequest) NewResponse(content string) *EventScancodePushResponse {

	return &EventScancodePushResponse{this.MessageHeader.getResponseHeader(), content}
}
