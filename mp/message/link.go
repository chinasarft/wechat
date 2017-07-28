package message

type LinkMessageHandler func(*LinkRequest) *LinkResponse

var linkMessageHandler LinkMessageHandler = nil

//Link 文本消息
type LinkRequest struct {
	MessageHeader
	Tile        string
	Description string
	Url         string
	MsgId       int64
}

type LinkResponse struct {
	MessageHeader
	Content string
}

func SetLinkMessageHandler(f LinkMessageHandler) {

	linkMessageHandler = f
}

//NewLink 初始化文本消息
func (this *LinkRequest) NewResponse(content string) *LinkResponse {

	return &LinkResponse{this.MessageHeader.getResponseHeader(), content}
}
