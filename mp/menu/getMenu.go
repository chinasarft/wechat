package menu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/token"
)

// TODO 可能对该struct需要添加一些方法
type RetrieveMenu struct {
	Normal     Menu    `json:"menu,omitempty"`
	Conditinal []*Menu `json:"conditionalmenu,omitempty"`
}

func GetMenu() (*RetrieveMenu, error) {
	accessToken := token.GetAccessToken()
	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=" + accessToken)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	menu := &RetrieveMenu{}
	err = json.Unmarshal(body, menu)
	if err != nil {
		return nil, err
	}

	return menu, nil
}
