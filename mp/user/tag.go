package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/token"
)

type Tag struct {
	Name  string `json:"name,omitempty"`
	Id    int    `json:"id,omitempty"`
	Count int    `json:"count,omitempty"`
}

func CreateTag(tagName string) (*Tag, error) {
	accessToken := token.GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=" + accessToken

	reqTag := &struct {
		Tag Tag `json:"tag"`
	}{}
	reqTag.Tag.Name = tagName

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqByte))
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

	result := &struct {
		Tag     *Tag   `json:"tag"`
		ErrCode int    `josn:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%d:%s", result.ErrCode, result.ErrMsg)
	}

	return result.Tag, nil
}

func UpdateTag(tagName string, id int) error {
	accessToken := token.GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/tags/update?access_token=" + accessToken

	reqTag := &struct {
		Tag Tag `json:"tag"`
	}{}
	reqTag.Tag.Name = tagName
	reqTag.Tag.Id = id

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqByte))
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

	result := &struct {
		ErrCode int    `josn:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("%d:%s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

func DeleteTag(id int) error {
	accessToken := token.GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=" + accessToken

	reqTag := &struct {
		Tag Tag `json:"tag"`
	}{}
	reqTag.Tag.Id = id

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqByte))
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

	result := &struct {
		ErrCode int    `josn:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("%d:%s", result.ErrCode, result.ErrMsg)
	}

	return nil
}

func GetTags() ([]*Tag, error) {
	accessToken := token.GetAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/tags/get?access_token=" + accessToken

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

	fmt.Println(string(body))
	list := &struct {
		Tags []*Tag `json:"tags"`
	}{}
	err = json.Unmarshal(body, list)
	if err != nil {
		return nil, err
	}

	return list.Tags, nil
}
