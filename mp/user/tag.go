package user

import (
	"encoding/json"
	"fmt"

	"github.com/chinasarft/wechat/mp/token"
)

type Tag struct {
	Name  string `json:"name,omitempty"`
	Id    int    `json:"id,omitempty"`
	Count int    `json:"count,omitempty"`
}

type TagList struct {
	Tags []*Tag `json:"tag"`
}

type ReqTag struct {
	Tag Tag `json:"tag"`
}

type TagWithError struct {
	Tag     *Tag   `json:"tag"`
	ErrCode int    `josn:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewReqTag(name string) *ReqTag {
	return &ReqTag{Tag{Name: name}}
}

func CreateTag(tagName string) (*Tag, error) {

	reqTag := NewReqTag(tagName)

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return nil, err
	}

	body, err := token.HttpPost("https://api.weixin.qq.com/cgi-bin/tags/create?access_token=",
		"application/json", reqByte)
	if err != nil {
		return nil, err
	}

	result := &TagWithError{}
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

	reqTag := &struct {
		Tag Tag `json:"tag"`
	}{}
	reqTag.Tag.Name = tagName
	reqTag.Tag.Id = id

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(reqByte))
	body, err := token.HttpPost("https://api.weixin.qq.com/cgi-bin/tags/update?access_token=",
		"application/json", reqByte)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
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

	reqTag := &struct {
		Tag Tag `json:"tag"`
	}{}
	reqTag.Tag.Id = id

	reqByte, err := json.MarshalIndent(reqTag, "", "\t")
	if err != nil {
		return err
	}

	body, err := token.HttpPost("https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=",
		"application/json", reqByte)
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

	body, err := token.HttpGet("https://api.weixin.qq.com/cgi-bin/tags/get?access_token=")

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
