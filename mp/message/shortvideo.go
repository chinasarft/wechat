package message

type ShortvideoMessageHandler func(*ShortvideoRequest) *ShortvideoResponse

var shortvideoMessageHandler ShortvideoMessageHandler = nil

//Shortvideo 文本消息
type ShortvideoRequest struct {
	MessageHeader
	MediaId      string
	ThumbMediaId string
	MsgId        int64
}

type ShortvideoResponse struct {
	MessageHeader
	Content string
}

func SetShortvideoMessageHandler(f ShortvideoMessageHandler) {

	shortvideoMessageHandler = f
}

//NewShortvideo 初始化文本消息
func (this *ShortvideoRequest) NewResponse(content string) *ShortvideoResponse {

	return &ShortvideoResponse{this.MessageHeader.getResponseHeader(), content}
}
