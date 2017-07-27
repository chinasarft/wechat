package token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	configFile = "config.json"
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
	App App      `json:"app"`
	Ak  WechatAk `json:"ak"`
}

var initOnce sync.Once
var refreshOnce sync.Once
var mutex sync.Mutex
var config Config

func Init() {
	initOnce.Do(loadConfig)
}

func loadConfig() {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	if time.Now().Unix() > config.Ak.ExpireTime {
		GetWechatAccessToken()
	}
	refreshOnce.Do(refreshToken)
}

func GetAccessToken() string {
	mutex.Lock()
	defer mutex.Unlock()
	return config.Ak.AccessToken
}

func refreshToken() {
	go func() {
		now := time.Now().Unix()
		timer := time.NewTimer(time.Second * time.Duration(config.Ak.ExpireTime-now))
		for {
			select {
			case <-timer.C:
				GetWechatAccessToken()
				timer.Reset(time.Second * time.Duration(config.Ak.ExpiresIn))
			}
		}
	}()
}

func GetWechatAccessToken() {
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
