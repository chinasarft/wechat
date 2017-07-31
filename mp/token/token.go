package token

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	configFile                    = "config.json"
	WECHAT_AES_KEY_LENGTH         = 43
	WECHAT_VALIDATE_TOKEN_MIN_LEN = 3
	WECHAT_VALIDATE_TOKEN_MAX_LEN = 32
)

type App struct {
	AppId     string `json:"appid"`
	AppSecret string `json:"secret"`
}

type WechatAk struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ExpireTime  int64  `json:expires_time`
}

type Config struct {
	App          App      `json:"app"`
	Ak           WechatAk `json:"ak"`
	Base64AesKey string   `json:"aes_key"`
	BinAesKey    []byte   `json:-`
	//https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319
	ValidateToken string `json:"validate_token"`
}

var initOnce sync.Once
var refreshOnce sync.Once
var mutex sync.Mutex
var config Config

func Init() {
	initOnce.Do(initConfig)
}

func initConfig() {

	loadConfigFile()

	checkAndDecodeAesKey()

	checkValidateToken()

	if isTokenExpire() {
		UPdateWechatAccessToken()
	}
	refreshOnce.Do(refreshToken)
}

func loadConfigFile() {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

func checkAndDecodeAesKey() {
	if len(config.Base64AesKey) != WECHAT_AES_KEY_LENGTH {
		panic("aes key wrong length")
	}
	aesKey, err := base64.StdEncoding.DecodeString(config.Base64AesKey + "=")
	if err != nil {
		panic(err)
	}
	config.BinAesKey = aesKey
}

func checkValidateToken() {
	validateTokenLen := len(config.ValidateToken)
	if validateTokenLen < WECHAT_VALIDATE_TOKEN_MIN_LEN ||
		validateTokenLen > WECHAT_VALIDATE_TOKEN_MAX_LEN {
		panic("validate token wrong length")
	}
}

func isTokenExpire() bool {
	return time.Now().Unix() > config.Ak.ExpireTime
}

func GetAccessToken() string {
	mutex.Lock()
	defer mutex.Unlock()
	return config.Ak.AccessToken
}

func GetAesKey() []byte {
	mutex.Lock()
	defer mutex.Unlock()
	return config.BinAesKey
}

func GetValidateToken() string {
	mutex.Lock()
	defer mutex.Unlock()
	return config.ValidateToken
}

func GetAppId() string {
	mutex.Lock()
	defer mutex.Unlock()
	return config.App.AppId
}

func refreshToken() {
	go func() {
		now := time.Now().Unix()
		timer := time.NewTimer(time.Second * time.Duration(config.Ak.ExpireTime-now))
		for {
			select {
			case <-timer.C:
				UPdateWechatAccessToken()
				timer.Reset(time.Second * time.Duration(config.Ak.ExpiresIn))
			}
		}
	}()
}

func UPdateWechatAccessToken() {
	mutex.Lock()
	defer mutex.Unlock()

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + config.App.AppId + "&secret=" + config.App.AppSecret

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ak := WechatAk{}
	err = json.Unmarshal(body, &ak)
	if err != nil {
		panic(err)
	}

	ak.ExpireTime = time.Now().Unix() + ak.ExpiresIn - 10
	config.Ak = ak

	writeToken()
}

func writeToken() {
	data, err := json.Marshal(&config)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(configFile, data, 0600); err != nil {
		panic(err)
	}
}
