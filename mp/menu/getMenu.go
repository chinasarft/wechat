package menu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/token"
)

type MixMenu struct {
	Buttons   []*Button  `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"` // 只在菜单查询接口返回的conditinalmenu包含这个字段
	MenuId    int64      `json:"menuid,omitempty"`    // 有个性化菜单时查询接口返回值包含这个字段
}

// TODO 可能对该struct需要添加一些方法
type RetrieveMenu struct {
	Normal     MixMenu    `json:"menu,omitempty"`
	Conditinal []*MixMenu `json:"conditionalmenu,omitempty"`
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
