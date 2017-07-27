// goproject project main.go
package main

import (
	"bytes"
	"fmt"

	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chinasarft/wechat/mp"
	"github.com/chinasarft/wechat/mp/menu"
	"github.com/chinasarft/wechat/mp/token"

	"github.com/gin-gonic/gin"
)

type WechatErr struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}
type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

func VlidateWechatServer(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	//验证微信连接

	if mp.ValidateWechatServer("6f68fe5452a9fee642d959410ab455af", timestamp, nonce, signature, echostr) {
		c.String(200, echostr)
	}
}

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func TextMsg(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// handle error
		fmt.Println("text msg fail")
	}
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	fmt.Println(requestBody)
	respContent, err := makeTextResponseBody(requestBody.ToUserName, requestBody.FromUserName, "welcome to yimi")
	if err == nil {
		fmt.Println("write response")
		c.Data(http.StatusOK, "text/xml", respContent)
	}
}

func main() {
	token.Init()
	engine := gin.New()
	engine.Static("/static", "static")
	weChatCoreGroupR := engine.Group("/wechat")
	{
		weChatCoreGroupR.GET("/connect", VlidateWechatServer)
		weChatCoreGroupR.POST("/connect", TextMsg)
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
