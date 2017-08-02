package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"
)

func DeleteAllMenu() error {
	accessToken := token.GetAccessToken()
	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=" + accessToken)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	requestResult := &message.RequstWechatResult{}
	err = json.Unmarshal(body, requestResult)
	if err != nil {
		return err
	}
	if requestResult.Errcode != 0 {
		return fmt.Errorf(string(body))
	}
	return nil
}

func DeleteConditionalMenu(menuId *MenuId) error {
	text, err := json.Marshal(menuId)
	if err != nil {
		return err
	}

	accessToken := token.GetAccessToken()
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token="+accessToken,
		"application/json", bytes.NewReader(text))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	requestResult := &message.RequstWechatResult{}
	err = json.Unmarshal(body, requestResult)
	if err != nil {
		return err
	}
	if requestResult.Errcode != 0 {
		return fmt.Errorf(string(body))
	}
	return nil
}
