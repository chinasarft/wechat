package message

type VoiceMessageHandler func(*VoiceRequest) *VoiceResponse

var voiceMessageHandler VoiceMessageHandler = nil

//Voice 文本消息
type VoiceRequest struct {
	MessageHeader
	MediaId     string
	Format      string
	Recognition string
	MsgId       int64
}

type VoiceResponse struct {
	MessageHeader
	Content string
}

func SetVoiceMessageHandler(f VoiceMessageHandler) {

	voiceMessageHandler = f
}

//NewVoice 初始化文本消息
func (this *VoiceRequest) NewResponse(content string) *VoiceResponse {

	return &VoiceResponse{this.MessageHeader.getResponseHeader(), content}
}
