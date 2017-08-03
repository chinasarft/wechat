// goproject project main.go
package main

import (
	"fmt"

	"net/http"
	"strconv"

	"github.com/chinasarft/wechat/mp/menu"
	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"

	"github.com/chinasarft/wechat/lib/flags"
	"github.com/gin-gonic/gin"
)

func ValidateWechatServer(c *gin.Context) {

	resp, err := message.ServerAuthentication(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Data(http.StatusOK, "", resp)
}

func textMsgHandler(r *message.TextRequest) *message.TextResponse {
	return r.NewResponse("text response:" + r.Content)
}
func locationMsgHandler(r *message.LocationRequest) *message.LocationResponse {
	return r.NewResponse("now you were stand at:" + r.Label)
}
func imageMsgHandler(r *message.ImageRequest) *message.ImageResponse {
	return r.NewResponse("your image at:" + r.PicUrl)
}
func voiceMsgHandler(r *message.VoiceRequest) *message.VoiceResponse {
	return r.NewResponse("your voice id:" + r.MediaId)
}
func videoMsgHandler(r *message.VideoRequest) *message.VideoResponse {
	return r.NewResponse("your vidoe id:" + r.MediaId + " " + r.ThumbMediaId)
}
func shortvideoMsgHandler(r *message.ShortvideoRequest) *message.ShortvideoResponse {
	return r.NewResponse("your shortvidoe id:" + r.MediaId + " " + r.ThumbMediaId)
}
func linkMsgHandler(r *message.LinkRequest) *message.LinkResponse {
	return r.NewResponse("your link tile:" + r.Tile + " " + r.Description)
}

func eventClickHandler(r *message.EventClickRequest) *message.EventClickResponse {
	return r.NewResponse("your key:" + r.EventKey)
}
func eventSubscribeHandler(r *message.EventSubscribeRequest) *message.EventSubscribeResponse {
	if r.EventKey == "" {
		return r.NewResponse("subscribe:" + string(r.Event))
	} else {
		return r.NewResponse("subscribe by scan:" + string(r.Event) + r.EventKey + r.Ticket)
	}
}
func eventUnsubscribeHandler(r *message.EventUnsubscribeRequest) *message.EventUnsubscribeResponse {
	return r.NewResponse("your evente:" + string(r.Event))
}
func eventViewHandler(r *message.EventViewRequest) *message.EventViewResponse {
	return r.NewResponse("your url:" + r.EventKey)
}
func eventLocationHandler(r *message.EventLocationRequest) *message.EventLocationResponse {
	return r.NewResponse("your eventlo:" + r.Latitude)
}
func eventLocationSelectHandler(r *message.EventLocationSelectRequest) *message.EventLocationSelectResponse {
	return r.NewResponse("your els:" + r.SendLocationInfo.Label + " " + r.SendLocationInfo.Poiname)
}
func eventScancodePushHandler(r *message.EventScancodePushRequest) *message.EventScancodePushResponse {
	return r.NewResponse("your esp:" + r.ScanCodeInfo.ScanType + " " + r.ScanCodeInfo.ScanResult)
}
func eventScancodeWaitmsgHandler(r *message.EventScancodeWaitmsgRequest) *message.EventScancodeWaitmsgResponse {
	return r.NewResponse("your esw:" + r.ScanCodeInfo.ScanType + " " + r.ScanCodeInfo.ScanResult)
}
func eventPicSysphotoHandler(r *message.EventPicSysphotoRequest) *message.EventPicSysphotoResponse {
	return r.NewResponse("your eps:" + strconv.Itoa(r.SendPicsInfo.Count) + " " + r.SendPicsInfo.PicList[0].PicMd5Sum)
}
func eventPicPhotoOrAlbumHandler(r *message.EventPicPhotoOrAlbumRequest) *message.EventPicPhotoOrAlbumResponse {
	return r.NewResponse("your epic_photo_or_album:" + strconv.Itoa(r.SendPicsInfo.Count) + " " + r.SendPicsInfo.PicList[0].PicMd5Sum)
}
func eventPicWeixinHandler(r *message.EventPicWeixinRequest) *message.EventPicWeixinResponse {
	return r.NewResponse("your epic_Weixin:" + strconv.Itoa(r.SendPicsInfo.Count) + " " + r.SendPicsInfo.PicList[0].PicMd5Sum)
}

func MessageGateway(c *gin.Context) {
	resp, err := message.HandleMessage(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Data(http.StatusOK, "", resp)
}

func main() {
	token.Init()
	fmt.Println(flags.GetTest())
	switch flags.GetTest() {
	case "freshMenu":
		freshMenu()
		m, err := menu.GetMenu()
		if err == nil {
			fmt.Printf("%+v\n", m)
		}
		return
	case "getMenu":
		m, err := menu.GetMenu()
		if err == nil {
			fmt.Printf("%+v\n", m)
		} else {
			fmt.Println("getMenu err:", err.Error())
		}
		return
	case "getToken":
		getToken()
		return
	case "cMenu":
		addConditionalMenu()
		return
	case "dcMenu":
		deleteConditionalMenu()
		return
	}

	message.SetTextMessageHandler(textMsgHandler)
	message.SetLocationMessageHandler(locationMsgHandler)
	message.SetImageMessageHandler(imageMsgHandler)
	message.SetVoiceMessageHandler(voiceMsgHandler)
	message.SetVideoMessageHandler(videoMsgHandler)
	message.SetShortvideoMessageHandler(shortvideoMsgHandler)
	message.SetLinkMessageHandler(linkMsgHandler)

	message.SetEventLocationMessageHandler(eventLocationHandler)
	message.SetEventClickMessageHandler(eventClickHandler)
	message.SetEventSubscribeMessageHandler(eventSubscribeHandler)
	message.SetEventUnsubscribeMessageHandler(eventUnsubscribeHandler)
	message.SetEventViewMessageHandler(eventViewHandler)
	message.SetEventLocationSelectMessageHandler(eventLocationSelectHandler)
	message.SetEventScancodePushMessageHandler(eventScancodePushHandler)
	message.SetEventScancodeWaitmsgMessageHandler(eventScancodeWaitmsgHandler)
	message.SetEventPicSysphotoMessageHandler(eventPicSysphotoHandler)
	message.SetEventPicPhotoOrAlbumMessageHandler(eventPicPhotoOrAlbumHandler)
	message.SetEventPicWeixinMessageHandler(eventPicWeixinHandler)

	engine := gin.New()
	engine.Static("/static", "static")
	weChatCoreGroupR := engine.Group("/wechat")
	{
		weChatCoreGroupR.GET("/connect", ValidateWechatServer)
		weChatCoreGroupR.POST("/connect", MessageGateway)
	}
	startServe(engine)
}

func deleteConditionalMenu() {
	menus, err := menu.GetMenu()
	if err != nil {
		fmt.Println(err)
		return
	}

	mids := menus.GetConditionalMenuId()
	l := len(mids)
	if l < 2 {
		fmt.Println("conditional menu less than 2. not delete")
		return
	}

	err = menu.DeleteConditionalMenu(mids[l-1])
	if err != nil {
		fmt.Println("delete conditional menu err:", err.Error())
		return
	}
	fmt.Println("delete conditional menu ok")
}

func getToken() {
	fmt.Println(http.StatusOK, "token:"+token.GetAccessToken()+"\n")
}

func freshMenu() {
	m := menu.NewMenu()

	clickButton := menu.NewClickButton("点击1", "key1")
	m.AddButton(clickButton)

	locationButton := menu.NewLocationSelectButton("位置2", "key2")
	m.AddButton(locationButton)

	menuButton := menu.NewMenuButton("菜单1")

	clickButton3_1 := menu.NewClickButton("点击3_1", "key3_1")
	menuButton.AddSubButton(clickButton3_1)
	scpButton3_2 := menu.NewScancodePushButton("扫码3_2", "key3_2")
	menuButton.AddSubButton(scpButton3_2)
	viewButton3_3 := menu.NewViewButton("跳转3_3", "http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/static/a.html")
	menuButton.AddSubButton(viewButton3_3)

	m.AddMenuButton(menuButton)

	err := m.Submit()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("freshmenu ok")
}

func addConditionalMenu() {
	m := menu.NewConditionalMenu()

	clickButton := menu.NewClickButton("点击c1", "keyc1")
	m.AddButton(clickButton)

	locationButton := menu.NewLocationSelectButton("位置c2", "keyc2")
	m.AddButton(locationButton)

	menuButton := menu.NewMenuButton("菜单c2")

	clickButton3_1 := menu.NewClickButton("点击c3_1", "keyc3_1")
	menuButton.AddSubButton(clickButton3_1)
	scpButton3_2 := menu.NewScancodePushButton("扫码c3_2", "keyc3_2")
	menuButton.AddSubButton(scpButton3_2)
	viewButton3_3 := menu.NewViewButton("跳转c3_3", "http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/static/a.html")
	menuButton.AddSubButton(viewButton3_3)

	m.AddMenuButton(menuButton)

	rule := menu.NewMatchRule()
	rule.SetSexToMale()
	m.AddMatchRule(rule)

	mid, err := m.Submit()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("freshmenu ok:", mid.MenuId)
}

func startServe(engine *gin.Engine) error {
	//if isHTTPs, certFilePath, keyFilePath := flags.DoHTTPs(); isHTTPs {
	//    return engine.RunTLS(flags.HostAndPort(), certFilePath, keyFilePath)
	//} else {
	return engine.Run(":80")
	//}
}
