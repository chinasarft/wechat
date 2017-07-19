// goproject project main.go
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"

	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WechatAk struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
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

//验证微信来源
func SignAccount(token, timestamp, nonce, signature, echostr string) bool {
	s := []string{token, timestamp, nonce}
	sort.Sort(sort.StringSlice(s)) //将token、timestamp、nonce三个参数进行字典序排序
	s0 := strings.Join(s, "")      //将三个参数字符串拼接成一个字符串
	t := sha1.New()                //sha1加密
	io.WriteString(t, s0)
	s1 := fmt.Sprintf("%x", t.Sum(nil))
	if signature == s1 { //与signature对比
		return true
	}
	return false
}

func GetWeChatCore(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	//验证微信连接
	if SignAccount("6f68fe5452a9fee642d959410ab455af", timestamp, nonce, signature, echostr) {
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

func GetWechatAccessToken() {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wxc013571e83b295ed&secret=2d30a2a5a835bd33825c44c1c4da468b"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	ak := WechatAk{}
	err = json.Unmarshal(body, &ak)
	if err != nil {
		fmt.Println(ak)
	}
}

func main() {
	engine := gin.New()
	weChatCoreGroupR := engine.Group("/wechat")
	{
		weChatCoreGroupR.GET("/connect", GetWeChatCore)
		weChatCoreGroupR.POST("/connect", TextMsg)
	}
	startServe(engine)

}

func startServe(engine *gin.Engine) error {
	//if isHTTPs, certFilePath, keyFilePath := flags.DoHTTPs(); isHTTPs {
	//    return engine.RunTLS(flags.HostAndPort(), certFilePath, keyFilePath)
	//} else {
	return engine.Run(":80")
	//}
}
