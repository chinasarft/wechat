package message

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"math/big"
	"sort"

	"github.com/chinasarft/wechat/mp/token"
)

// MsgType 基本消息类型
type MsgType string

// EventType 事件类型
type EventType string

const (
	//MsgTypeText 表示文本消息
	MsgTypeText MsgType = "text"
	//MsgTypeImage 表示图片消息
	MsgTypeImage = "image"
	//MsgTypeVoice 表示语音消息
	MsgTypeVoice = "voice"
	//MsgTypeVideo 表示视频消息
	MsgTypeVideo = "video"
	//MsgTypeShortVideo 表示短视频消息[限接收]
	MsgTypeShortVideo = "shortvideo"
	//MsgTypeLocation 表示坐标消息[限接收]
	MsgTypeLocation = "location"
	//MsgTypeLink 表示链接消息[限接收]
	MsgTypeLink = "link"
	//MsgTypeMusic 表示音乐消息[限回复]
	MsgTypeMusic = "music"
	//MsgTypeNews 表示图文消息[限回复]
	MsgTypeNews = "news"
	//MsgTypeTransfer 表示消息消息转发到客服
	MsgTypeTransfer = "transfer_customer_service"
	//MsgTypeEvent 表示事件推送消息
	MsgTypeEvent = "event"
)

const (
	//EventSubscribe 订阅
	EventSubscribe EventType = "subscribe"
	//EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	//EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	//EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	//EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	//EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	//EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	//EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	//EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	//EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	//EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	//EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
)

//MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	MessageHeader

	//基本消息
	MsgId        int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	PicUrl       string  `xml:"PicUrl"`
	MediaId      string  `xml:"MediaId"`
	Format       string  `xml:"Format"`
	ThumbMediaId string  `xml:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64 `xml:"Scale"`
	Label        string  `xml:"Label"`
	Title        string  `xml:"Title"`
	Description  string  `xml:"Description"`
	Url          string  `xml:"Url"`
	Recognition  string  `xml:"Recognition"`

	//事件相关
	Event            EventType        `xml:"Event"`
	EventKey         string           `xml:"EventKey"`
	Ticket           string           `xml:"Ticket"`
	Latitude         string           `xml:"Latitude"`
	Longitude        string           `xml:"Longitude"`
	Precision        string           `xml:"Precision"`
	MenuId           string           `xml:"MenuId"`
	ScanCodeInfo     ScanCodeInfo     `xml:"ScanCodeInfo"`
	SendPicsInfo     SendPicsInfo     `xml:"SendPicsInfo"`
	SendLocationInfo SendLocationInfo `xml:"SendLocationInfo"`
}

type ScanCodeInfo struct {
	ScanType   string `xml:"ScanType"`
	ScanResult string `xml:"ScanResult"`
}

type SendPicsInfo struct {
	Count   int32      `xml:"Count"`
	PicList []EventPic `xml:"PicList>item"`
}

type SendLocationInfo struct {
	LocationX float64 `xml:"Location_X"`
	LocationY float64 `xml:"Location_Y"`
	Scale     float64 `xml:"Scale"`
	Label     string  `xml:"Label"`
	Poiname   string  `xml:"Poiname"`
}

//EventPic 发图事件推送
type EventPic struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

//EncryptedXMLMsg 安全模式下的消息体
type EncryptMsg struct {
	XMLName    xml.Name `xml:"xml" json:"-"`
	ToUserName string   `xml:"ToUserName" json:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"    json:"Encrypt"`
}

//ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	Encrypt      CDATA    `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature CDATA    `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        CDATA    `xml:"Nonce"        json:"Nonce"`
}

// CommonToken 消息中通用的结构
type MessageHeader struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      MsgType  `xml:"MsgType"`
}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

func (h MessageHeader) getResponseHeader() MessageHeader {
	respH := h
	respH.FromUserName = h.ToUserName
	respH.ToUserName = h.FromUserName
	respH.CreateTime = time.Now().Unix()
	// TODO all message return MsgTypeText is OK?
	respH.MsgType = MsgTypeText
	return respH
}

func upperFirstChar(s string) string {
	ls := []byte(strings.ToLower(s))
	return strings.ToUpper(string(ls[0:1])) + string(ls[1:])
}

func (this *MixMessage) getMethodName() (m []string) {

	if this.MsgType != MsgTypeEvent {
		msgType := upperFirstChar(string(this.MsgType))
		m = append(m, "Get"+msgType+"Request")
		m = append(m, msgType+"RequestHandler")
		return
	}
	event := upperFirstChar(string(this.Event))
	m = append(m, "GetEvent"+event+"Request")
	m = append(m, "Event"+event+"RequestHandler")
	return
}

func (this *MixMessage) GetTextRequest() *TextRequest {
	textRequest := &TextRequest{this.MessageHeader, this.Content, this.MsgId}
	return textRequest
}

func (this *MixMessage) TextRequestHandler(textRequestI interface{}) *TextResponse {
	textRequest := textRequestI.(*TextRequest)
	if textMessageHandler != nil {
		return textMessageHandler(textRequest)
	}
	return &TextResponse{}
}

func (this *MixMessage) GetImageRequest() *ImageRequest {
	imageRequest := &ImageRequest{this.MessageHeader, this.PicUrl, this.MediaId, this.MsgId}
	return imageRequest
}

func (this *MixMessage) ImageRequestHandler(imageRequestI interface{}) *ImageResponse {
	imageRequest := imageRequestI.(*ImageRequest)
	if imageMessageHandler != nil {
		return imageMessageHandler(imageRequest)
	}
	return &ImageResponse{}
}

func (this *MixMessage) GetVoiceRequest() *VoiceRequest {
	voiceRequest := &VoiceRequest{this.MessageHeader, this.MediaId, this.Format, this.Recognition, this.MsgId}
	return voiceRequest
}

func (this *MixMessage) VoiceRequestHandler(voiceRequestI interface{}) *VoiceResponse {
	voiceRequest := voiceRequestI.(*VoiceRequest)
	if voiceMessageHandler != nil {
		return voiceMessageHandler(voiceRequest)
	}
	return &VoiceResponse{}
}

func (this *MixMessage) GetVideoRequest() *VideoRequest {
	videoRequest := &VideoRequest{this.MessageHeader, this.MediaId, this.ThumbMediaId, this.MsgId}
	return videoRequest
}

func (this *MixMessage) VideoRequestHandler(videoRequestI interface{}) *VideoResponse {
	videoRequest := videoRequestI.(*VideoRequest)
	if videoMessageHandler != nil {
		return videoMessageHandler(videoRequest)
	}
	return &VideoResponse{}
}

func (this *MixMessage) GetShortvideoRequest() *ShortvideoRequest {
	shortvideoRequest := &ShortvideoRequest{this.MessageHeader, this.MediaId, this.ThumbMediaId, this.MsgId}
	return shortvideoRequest
}

func (this *MixMessage) ShortvideoRequestHandler(shortvideoRequestI interface{}) *ShortvideoResponse {
	shortvideoRequest := shortvideoRequestI.(*ShortvideoRequest)
	if shortvideoMessageHandler != nil {
		return shortvideoMessageHandler(shortvideoRequest)
	}
	return &ShortvideoResponse{}
}

func (this *MixMessage) GetLinkRequest() *LinkRequest {
	linkRequest := &LinkRequest{this.MessageHeader, this.Title, this.Description, this.Url, this.MsgId}
	return linkRequest
}

func (this *MixMessage) LinkRequestHandler(linkRequestI interface{}) *LinkResponse {
	linkRequest := linkRequestI.(*LinkRequest)
	if linkMessageHandler != nil {
		return linkMessageHandler(linkRequest)
	}
	return &LinkResponse{}
}

func (this *MixMessage) GetEventClickRequest() *EventClickRequest {
	eventClickRequest := &EventClickRequest{this.MessageHeader, this.Event, this.EventKey}
	return eventClickRequest
}

func (this *MixMessage) EventClickRequestHandler(eventClickRequestI interface{}) *EventClickResponse {
	eventClickRequest := eventClickRequestI.(*EventClickRequest)
	if eventClickMessageHandler != nil {
		return eventClickMessageHandler(eventClickRequest)
	}
	return &EventClickResponse{}
}

func (this *MixMessage) GetEventViewRequest() *EventViewRequest {
	viewRequest := &EventViewRequest{this.MessageHeader, this.Event, this.EventKey, this.MenuId}
	return viewRequest
}

func (this *MixMessage) EventViewRequestHandler(eventViewRequestI interface{}) *EventViewResponse {
	eventViewRequest := eventViewRequestI.(*EventViewRequest)
	if eventViewMessageHandler != nil {
		return eventViewMessageHandler(eventViewRequest)
	}
	return &EventViewResponse{}
}

func (this *MixMessage) GetEventLocationRequest() *EventLocationRequest {
	viewRequest := &EventLocationRequest{this.MessageHeader, this.Event, this.Latitude, this.Longitude, this.Precision}
	return viewRequest
}

func (this *MixMessage) EventLocationRequestHandler(eventLocationRequestI interface{}) *EventLocationResponse {
	eventLocationRequest := eventLocationRequestI.(*EventLocationRequest)
	if eventLocationMessageHandler != nil {
		return eventLocationMessageHandler(eventLocationRequest)
	}
	return &EventLocationResponse{}
}

func (this *MixMessage) GetEventLocation_selectRequest() *EventLocationSelectRequest {
	eventLocationSelectRequest := &EventLocationSelectRequest{this.MessageHeader, this.Event, this.EventKey, this.SendLocationInfo}
	return eventLocationSelectRequest
}

func (this *MixMessage) EventLocation_selectRequestHandler(eventLocationSelectRequestI interface{}) *EventLocationSelectResponse {
	eventLocationSelectRequest := eventLocationSelectRequestI.(*EventLocationSelectRequest)
	if eventLocationSelectMessageHandler != nil {
		return eventLocationSelectMessageHandler(eventLocationSelectRequest)
	}
	return &EventLocationSelectResponse{}
}

func (this *MixMessage) GetLocationRequest() *LocationRequest {
	locationRequest := &LocationRequest{this.MessageHeader, this.LocationX, this.LocationY, this.Scale, this.Label, this.MsgId}
	return locationRequest
}

func (this *MixMessage) LocationRequestHandler(locationRequestI interface{}) *LocationResponse {
	locationRequest := locationRequestI.(*LocationRequest)
	if locationMessageHandler != nil {
		return locationMessageHandler(locationRequest)
	}
	return &LocationResponse{}
}

func (this *MessageHeader) SetToUserName(toUserName string) {
	this.ToUserName = toUserName
}

func (this *MessageHeader) SetFromUserName(fromUserName string) {
	this.FromUserName = fromUserName
}

func (this *MessageHeader) SetCreateTime(createTime int64) {
	this.CreateTime = createTime
}

func (this *MessageHeader) SetMsgType(msgType MsgType) {
	this.MsgType = msgType
}

func getMixMessage(rawMsg []byte) (*MixMessage, error) {

	mixMsg := &MixMessage{}
	err := xml.Unmarshal(rawMsg, mixMsg)
	if err != nil {
		return nil, err
	}
	return mixMsg, nil
}

func handlePlain(mixMsg *MixMessage) interface{} {

	mixMsgValue := reflect.ValueOf(mixMsg)
	methodName := mixMsg.getMethodName()

	request := mixMsgValue.MethodByName(methodName[0]).Call(nil)

	params := make([]reflect.Value, 1)
	params[0] = request[0]
	response := mixMsgValue.MethodByName(methodName[1]).Call(params)

	return response[0].Interface()
}

func aesDecrypt(base64Data string, aesKey []byte) ([]byte, error) {
	cipherData, _ := base64.StdEncoding.DecodeString(base64Data)
	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, fmt.Errorf("crypto/cipher: ciphertext size is not multiple of aes key length")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	//signature := strings.Join(r.Form["signature"], "")
	encryptType := strings.Join(r.Form["encrypt_type"], "")
	msgSignature := strings.Join(r.Form["msg_signature"], "")

	var rawMsg []byte
	var err error
	var isSafeModel bool = false

	if r.Method == http.MethodPost {
		if encryptType == "aes" {
			isSafeModel = true
			fmt.Println("Wechat Service: in safe mode")
			encryptMsg := parseEncryptRequestBody(r)

			//Validate msg signature
			if !validateMsg(timestamp, nonce, encryptMsg.Encrypt, msgSignature) {
				fmt.Println("Wechat Service: msg_signature is invalid")
				return
			}
			fmt.Println("Wechat Service: msg_signature validation is ok!")

			fmt.Println("cipher:", encryptMsg.Encrypt)

			aesKey := token.GetAesKey()
			plainData, err := aesDecrypt(encryptMsg.Encrypt, aesKey)
			if err != nil {
				fmt.Println(len(aesKey), len(encryptMsg.Encrypt), "11:", err)
				return
			}
			fmt.Println("plain:", string(plainData))

			rawMsg, err = checkAndGetPlain(plainData)
			if err != nil {
				fmt.Println("checkAndGetPlain:", err)
				return
			}
		} else {
			rawMsg, err = ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("read err:", err)
				return
			}

		}

	} else if r.Method == http.MethodGet {
		fmt.Println("get method")
		return
	} else {
		fmt.Println(r.Method, "method")
		return
	}

	fmt.Println("rawXml:", string(rawMsg))

	mixMsg, err := getMixMessage(rawMsg)
	if err != nil {
		fmt.Println("1:", err)
		return
	}
	result := handlePlain(mixMsg)

	respByte, err := xml.MarshalIndent(result, "", "")
	if err != nil {
		fmt.Println("2:", err)
		return
	}
	fmt.Println("plainData:", string(respByte))
	if isSafeModel {
		cipherRespByte, err := makeEncryptResponseXml(respByte)
		if err != nil {
			fmt.Printf("makeEncryptResponseXml:", err)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		fmt.Println("response:", string(cipherRespByte))
		w.Write(cipherRespByte)
	} else {
		w.Write(respByte)
	}
	return
}

func parseEncryptRequestBody(r *http.Request) *EncryptMsg {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("9:", err)
		return nil
	}
	encryptMsg := &EncryptMsg{}
	xml.Unmarshal(body, encryptMsg)
	return encryptMsg
}

func validateMsg(timestamp, nonce, msgEncrypt, msgSignatureIn string) bool {
	msgSignatureGen := makeMsgSignature(timestamp, nonce, msgEncrypt)
	if msgSignatureGen != msgSignatureIn {
		return false
	}
	return true
}

func makeMsgSignature(timestamp, nonce, msg_encrypt string) string {
	validateToken := token.GetValidateToken()
	sl := []string{validateToken, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

//http://studygolang.com/articles/2213

func checkAndGetPlain(plainText []byte) ([]byte, error) {

	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	fmt.Println("length:", string(plainText[20:20+length]))

	appId := token.GetAppId()
	// appID validation
	appIdstart := 20 + length
	id := plainText[appIdstart : int(appIdstart)+len(appId)]
	if string(id) != appId {
		fmt.Println("Wechat Service: appid is invalid!")
		return nil, fmt.Errorf("Appid is invalid")
	}
	fmt.Println("Wechat Service: appid validation is ok!")

	return plainText[20 : 20+length], nil

}

func makeEncryptResponseXml(plainData []byte) ([]byte, error) {

	// Encrypt part2: Length bytes
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(plainData)))
	if err != nil {
		fmt.Println("Binary write err:", err)
		return nil, err
	}
	bodyLength := buf.Bytes()

	// Encrypt part1: Random bytes
	randomBytes := []byte("abcdefghijklmnop")

	appId := token.GetAppId()
	// Encrypt Part, with part4 - appID
	decoratePlainData := bytes.Join([][]byte{randomBytes, bodyLength, plainData, []byte(appId)}, nil)
	fmt.Println("decoratePlainData:", string(decoratePlainData[0:16]),
		decoratePlainData[16:20], string(decoratePlainData[20:]))
	cipherData, err := aesEncrypt(decoratePlainData, token.GetAesKey())
	if err != nil {
		return nil, fmt.Errorf("aesEncrypt error")
	}

	responseEncryptMsg := &ResponseEncryptMsg{}
	responseEncryptMsg.Encrypt = CDATA(base64.StdEncoding.EncodeToString(cipherData))

	numNonce, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	strNonce := numNonce.String()
	responseEncryptMsg.Nonce = CDATA(strNonce)

	responseEncryptMsg.Timestamp = time.Now().Unix()
	sig := makeMsgSignature(
		strconv.Itoa(int(responseEncryptMsg.Timestamp)),
		string(responseEncryptMsg.Nonce), string(responseEncryptMsg.Encrypt))
	responseEncryptMsg.MsgSignature = CDATA(sig)

	return xml.MarshalIndent(responseEncryptMsg, "", "")
}

func aesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey)
	plainData = PKCS7Padding(plainData, k)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)

	return cipherData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
