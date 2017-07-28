package message

type LocationMessageHandler func(*LocationRequest) *LocationResponse

var locationMessageHandler LocationMessageHandler = nil

//Location 文本消息
type LocationRequest struct {
	MessageHeader
	LocationX float64
	LocationY float64
	Scale     float64
	Label     string
	MsgId     int64
}

type LocationResponse struct {
	MessageHeader
	Content string
}

func SetLocationMessageHandler(f LocationMessageHandler) {

	locationMessageHandler = f
}

//NewLocation 初始化文本消息
func (this *LocationRequest) NewResponse(content string) *LocationResponse {

	return &LocationResponse{this.MessageHeader.getResponseHeader(), content}
}
