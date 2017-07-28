package message

type ImageMessageHandler func(*ImageRequest) *ImageResponse

var imageMessageHandler ImageMessageHandler = nil

//Image 文本消息
type ImageRequest struct {
	MessageHeader
	PicUrl  string
	MediaId string
	MsgId   int64
}

type ImageResponse struct {
	MessageHeader
	Content string
}

func SetImageMessageHandler(f ImageMessageHandler) {

	imageMessageHandler = f
}

//NewImage 初始化文本消息
func (this *ImageRequest) NewResponse(content string) *ImageResponse {

	return &ImageResponse{this.MessageHeader.getResponseHeader(), content}
}
