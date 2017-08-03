package message

type EventPicSysphotoMessageHandler func(*EventPicSysphotoRequest) *EventPicSysphotoResponse

var eventPicSysphotoMessageHandler EventPicSysphotoMessageHandler = nil

//EventPicSysphoto 文本消息
type EventPicSysphotoRequest struct {
	MessageHeader
	Event        EventType
	EventKey     string
	SendPicsInfo SendPicsInfo
}

type EventPicSysphotoResponse struct {
	MessageHeader
	Content string
}

func SetEventPicSysphotoMessageHandler(f EventPicSysphotoMessageHandler) {

	eventPicSysphotoMessageHandler = f
}

//NewEventPicSysphoto 初始化文本消息
func (this *EventPicSysphotoRequest) NewResponse(content string) *EventPicSysphotoResponse {

	return &EventPicSysphotoResponse{this.MessageHeader.getResponseHeader(), content}
}
