package message

type EventPicWeixinMessageHandler func(*EventPicWeixinRequest) *EventPicWeixinResponse

var eventPicWeixinMessageHandler EventPicWeixinMessageHandler = nil

//EventPicWeixin 文本消息
type EventPicWeixinRequest struct {
	MessageHeader
	Event        EventType
	EventKey     string
	SendPicsInfo SendPicsInfo
}

type EventPicWeixinResponse struct {
	MessageHeader
	Content string
}

func SetEventPicWeixinMessageHandler(f EventPicWeixinMessageHandler) {

	eventPicWeixinMessageHandler = f
}

//NewEventPicWeixin 初始化文本消息
func (this *EventPicWeixinRequest) NewResponse(content string) *EventPicWeixinResponse {

	return &EventPicWeixinResponse{this.MessageHeader.getResponseHeader(), content}
}
