package message

type TextMessageHandler func(*TextRequest) *TextResponse

var textMessageHandler TextMessageHandler = nil

//Text 文本消息
type TextRequest struct {
	MessageHeader
	Content string
	MsgId   int64
}

type TextResponse struct {
	MessageHeader
	Content string
}

func SetTextMessageHandler(f TextMessageHandler) {

	textMessageHandler = f
}

//NewText 初始化文本消息
func (this *TextRequest) NewResponse(content string) *TextResponse {

	return &TextResponse{this.MessageHeader.getResponseHeader(), content}
}
