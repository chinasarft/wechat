package menu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"
)

func DeleteMenu() error {
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
