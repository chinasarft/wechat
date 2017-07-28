package message

type VideoMessageHandler func(*VideoRequest) *VideoResponse

var videoMessageHandler VideoMessageHandler = nil

//Video 文本消息
type VideoRequest struct {
	MessageHeader
	MediaId      string
	ThumbMediaId string
	MsgId        int64
}

type VideoResponse struct {
	MessageHeader
	Content string
}

func SetVideoMessageHandler(f VideoMessageHandler) {

	videoMessageHandler = f
}

//NewVideo 初始化文本消息
func (this *VideoRequest) NewResponse(content string) *VideoResponse {

	return &VideoResponse{this.MessageHeader.getResponseHeader(), content}
}
