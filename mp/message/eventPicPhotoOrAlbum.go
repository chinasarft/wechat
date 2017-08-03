package message

type EventPicPhotoOrAlbumMessageHandler func(*EventPicPhotoOrAlbumRequest) *EventPicPhotoOrAlbumResponse

var eventPicPhotoOrAlbumMessageHandler EventPicPhotoOrAlbumMessageHandler = nil

//EventPicPhotoOrAlbum 文本消息
type EventPicPhotoOrAlbumRequest struct {
	MessageHeader
	Event        EventType
	EventKey     string
	SendPicsInfo SendPicsInfo
}

type EventPicPhotoOrAlbumResponse struct {
	MessageHeader
	Content string
}

func SetEventPicPhotoOrAlbumMessageHandler(f EventPicPhotoOrAlbumMessageHandler) {

	eventPicPhotoOrAlbumMessageHandler = f
}

//NewEventPicPhotoOrAlbum 初始化文本消息
func (this *EventPicPhotoOrAlbumRequest) NewResponse(content string) *EventPicPhotoOrAlbumResponse {

	return &EventPicPhotoOrAlbumResponse{this.MessageHeader.getResponseHeader(), content}
}
