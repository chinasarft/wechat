// goproject project main.go
package main

import (
	"fmt"

	"net/http"
	"net/url"
	"strconv"

	"github.com/chinasarft/wechat/mp/menu"
	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"
	"github.com/chinasarft/wechat/mp/user"

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
func eventScanHandler(r *message.EventScanRequest) *message.EventScanResponse {
	return r.NewResponse("event scan:" + string(r.Event) + r.EventKey + r.Ticket)
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

func oauthRedirect(c *gin.Context) {
	fmt.Println("oauth redirect")
	redirectUrl := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + token.GetAppId()
	redirectUrl += "&redirect_uri="
	redirectUrl += url.QueryEscape("http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/test/oauth/ok")
	redirectUrl += "&response_type=code"
	redirectUrl += "&scope=snsapi_base"
	redirectUrl += "&state=s123#wechat_redirect"
	c.Redirect(http.StatusFound, redirectUrl)
	//"https://ohorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE")
}

func oauthRedirect2(c *gin.Context) {
	fmt.Println("oauth redirect2")
	oauthUrl := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + token.GetAppId()
	oauthUrl += "&redirect_uri="
	oauthUrl += url.QueryEscape("http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/test/oauth/ok")
	oauthUrl += "&response_type=code"
	oauthUrl += "&scope=snsapi_userinfo"
	oauthUrl += "&state=s123#wechat_redirect"

	pre := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width,initial-scale=1,user-scalable=0">
    <title>WeUI</title>
    <link rel="stylesheet" href="https://cdn.bootcss.com/weui/1.1.2/style/weui.min.css">
</head>
<body ontouchstart>
  <div>
  <p>
  获取你的公开信息(昵称等)
  </p>
  </div>
  <a href="`

	last := `" class="weui-btn weui-btn_primary">授权</a>
</body>
</html>`
	c.Data(http.StatusOK, "text/html", []byte(pre+oauthUrl+last))
	//c.String(http.StatusOK, pre+oauthUrl+last)
	//"https://ohorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE")
}

func oauthOk(c *gin.Context) {
	if code := c.Query("code"); code != "" {
		fmt.Println("code:----->", code)
	} else {
		fmt.Println("----> no code")
	}
	c.Data(http.StatusOK, "", nil)
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
	case "getUsers":
		getSubscribeUserList()
		return
	case "getTags":
		getTagsList()
		return
	case "addTags":
		addTagsList()
		return
	case "upTags":
		updateTag()
		return
	case "delTags":
		deleteTag()
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

	message.SetEventScanMessageHandler(eventScanHandler)
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
	testGroupR := engine.Group("/test")
	{
		testGroupR.GET("/redirect/oauth", oauthRedirect)
		testGroupR.GET("/redirect/oauth2", oauthRedirect2)
		testGroupR.GET("/oauth/ok", oauthOk)
	}
	weChatCoreGroupR := engine.Group("/wechat")
	{
		weChatCoreGroupR.GET("/connect", ValidateWechatServer)
		weChatCoreGroupR.POST("/connect", MessageGateway)
	}
	startServe(engine)
}

func deleteTag() {
	err := user.DeleteTag(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("delete tag ok")
}

func updateTag() {
	err := user.UpdateTag("测试tag1", 101)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("update tag ok")
}

func addTagsList() {
	tag, err := user.CreateTag("测试tag")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tag)
}

func getTagsList() {
	tags, err := user.GetTags()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(tags); i++ {
		fmt.Println(tags[i])
	}
}

func getSubscribeUserList() {
	users, err := user.GetUserList("")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(users.Data.OpenIdList); i++ {
		fmt.Println(users.Data.OpenIdList[0], users)
	}

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
	viewButton3_4 := menu.NewViewButton("跳转3_4", "http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/test/redirect/oauth")
	menuButton.AddSubButton(viewButton3_4)
	viewButton3_5 := menu.NewViewButton("跳转3_5", "http://ec2-13-112-47-201.ap-northeast-1.compute.amazonaws.com/test/redirect/oauth2")
	menuButton.AddSubButton(viewButton3_5)

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
