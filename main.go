// goproject project main.go
package main

import (
	"bytes"
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp"
	"github.com/chinasarft/wechat/mp/menu"
	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"

	"github.com/gin-gonic/gin"
)

type WechatErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func ValidateWechatServer(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	//验证微信连接

	if mp.ValidateWechatServer(token.GetValidateToken(), timestamp, nonce, signature, echostr) {
		c.String(200, echostr)
	}
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
func eventViewHandler(r *message.EventViewRequest) *message.EventViewResponse {
	return r.NewResponse("your url:" + r.EventKey)
}
func eventLocationHandler(r *message.EventLocationRequest) *message.EventLocationResponse {
	return r.NewResponse("your eventlo:" + r.Latitude)
}
func eventLocationSelectHandler(r *message.EventLocationSelectRequest) *message.EventLocationSelectResponse {
	return r.NewResponse("your els:" + r.SendLocationInfo.Label + " " + r.SendLocationInfo.Poiname)
}

func MessageGateway(c *gin.Context) {
	message.Handle(c.Writer, c.Request)
}

func main() {
	token.Init()
	message.SetTextMessageHandler(textMsgHandler)
	message.SetLocationMessageHandler(locationMsgHandler)
	message.SetImageMessageHandler(imageMsgHandler)
	message.SetVoiceMessageHandler(voiceMsgHandler)
	message.SetVideoMessageHandler(videoMsgHandler)
	message.SetShortvideoMessageHandler(shortvideoMsgHandler)
	message.SetLinkMessageHandler(linkMsgHandler)

	message.SetEventLocationMessageHandler(eventLocationHandler)
	message.SetEventClickMessageHandler(eventClickHandler)
	message.SetEventViewMessageHandler(eventViewHandler)
	message.SetEventLocationSelectMessageHandler(eventLocationSelectHandler)

	engine := gin.New()
	engine.Static("/static", "static")
	weChatCoreGroupR := engine.Group("/wechat")
	{
		weChatCoreGroupR.GET("/connect", ValidateWechatServer)
		weChatCoreGroupR.POST("/connect", MessageGateway)
	}

	mygroup := engine.Group("/test")
	{
		mygroup.POST("/menu", freshMenu)
		mygroup.GET("/token", getToken)
	}
	startServe(engine)
}

func getToken(c *gin.Context) {
	c.String(http.StatusOK, "token:"+token.GetAccessToken()+"\n")
}
func freshMenu(c *gin.Context) {
	m := menu.NewMenu()

	clickButton := menu.NewClickButton("点击1", "key1")
	m.AddButton(clickButton)

	locationButton := menu.NewLocationSelectButton("位置2", "key2")
	m.AddButton(locationButton)

	menuButton := menu.NewMenuButton("菜单")

	clickButton3_1 := menu.NewClickButton("点击3_1", "key3_1")
	menuButton.AddSubButton(clickButton3_1)
	scpButton3_2 := menu.NewScancodePushButton("扫码3_2", "key3_2")
	menuButton.AddSubButton(scpButton3_2)
	viewButton3_3 := menu.NewViewButton("跳转3_3", "www.bing.com")
	menuButton.AddSubButton(viewButton3_3)

	m.AddMenuButton(menuButton)

	text, err := m.GetJsonByte()
	if err != nil {
		fmt.Println("marshal error")
		return
	}

	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/menu/create?access_token="+token.GetAccessToken(),
		"application/json", bytes.NewReader(text))
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(string(body))

}

func startServe(engine *gin.Engine) error {
	//if isHTTPs, certFilePath, keyFilePath := flags.DoHTTPs(); isHTTPs {
	//    return engine.RunTLS(flags.HostAndPort(), certFilePath, keyFilePath)
	//} else {
	return engine.Run(":80")
	//}
}
