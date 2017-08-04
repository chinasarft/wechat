package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/token"
)

type UserList struct {
	Totoal int64 `json:"total"`
	Count  int64 `json:"count"`
	Data   struct {
		OpenIdList []string `json:"openid,omitempty"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
	ErrCode    int    `josn:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

//{"total":2,"count":2,"data":{"openid":["","OPENID1","OPENID2"]},"next_openid":"NEXT_OPENID"}

func GetUserList(nextOpenid string) (*UserList, error) {
	accessToken := token.GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + accessToken
	if nextOpenid != "" {
		url = url + "&next_openid=" + nextOpenid
	}

	resp, err := http.Get(url)
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

	list := &UserList{}
	err = json.Unmarshal(body, list)
	if err != nil {
		return nil, err
	}
	if list.ErrCode != 0 {
		return nil, fmt.Errorf("%d:%s", list.ErrCode, list.ErrMsg)
	}

	return list, nil
}
