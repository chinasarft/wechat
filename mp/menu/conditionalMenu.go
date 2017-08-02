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

const (
	SexMale                   = "1"
	SexFemale                 = "2"
	ClientPlatformTypeAndroid = "1"
	ClientPlatformTypeIOS     = "2"
	ClientPlatformTypeOthers  = "3"
)

type MenuId struct {
	MenuId int64 `json:"menuid,omitempty"`
}

type ConditionalMenu struct {
	Buttons   []*Button  `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
}

type MatchRule struct {
	GroupId            *int64 `json:"group_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
	TagId              string `json:"tag_id,omitempty"`
}

func NewConditionalMenu() *ConditionalMenu {
	return &ConditionalMenu{}
}

func (this *ConditionalMenu) AddButton(button *Button) error {
	if len(this.Buttons) >= MAX_FIRST_CLASS_MENU {
		return ErrMaxMenu
	}
	this.Buttons = append(this.Buttons, button)
	return nil
}

func (this *ConditionalMenu) AddMenuButton(menuButton *MenuButton) error {
	if len(this.Buttons) >= MAX_FIRST_CLASS_MENU {
		return ErrMaxMenu
	}
	button := (*Button)(menuButton)
	this.Buttons = append(this.Buttons, button)
	return nil
}

func (this *ConditionalMenu) GetJsonByte() ([]byte, error) {
	return json.MarshalIndent(this, "", "\t")
}

func (this *ConditionalMenu) AddMatchRule(matchRule *MatchRule) {
	this.MatchRule = matchRule
}

func (this *ConditionalMenu) Submit() (*MenuId, error) {

	text, err := this.GetJsonByte()
	if err != nil {
		return nil, err
	}

	accessToken := token.GetAccessToken()
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token="+accessToken,
		"application/json", bytes.NewReader(text))
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

	// two kind json respose
	// {"errcode":65304,"errmsg":"match rule empty hint: [d7XCTa0703vr18]"}
	// {"menuid":459976677}
	result := &message.RequstWechatResult{}
	if err = json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	menuId := &MenuId{}
	if result.Errcode == 0 && result.Errmsg == "" {
		err = json.Unmarshal(body, menuId)
		if err != nil {
			return nil, err
		}
		return menuId, nil
	}

	return nil, fmt.Errorf(string(body))
}

func NewMatchRule() *MatchRule {
	return &MatchRule{}
}

func (this *MatchRule) SetSexToMale() {
	this.Sex = SexMale
}

func (this *MatchRule) SetSexToFemale() {
	this.Sex = SexFemale
}

func (this *MatchRule) SetPlatformToAndroid() {
	this.ClientPlatformType = ClientPlatformTypeAndroid
}

func (this *MatchRule) SetPlatformToIOS() {
	this.ClientPlatformType = ClientPlatformTypeIOS
}

func (this *MatchRule) SetPlatformToOthers() {
	this.ClientPlatformType = ClientPlatformTypeOthers
}

func (this *MatchRule) SetGroudId(id int64) {
	gid := new(int64)
	*gid = id
	this.GroupId = gid
}

func (this *MatchRule) SetCountry(country string) {
	this.Country = country
}

func (this *MatchRule) SetCountryProvince(country, province string) {
	this.Country = country
	this.Province = province
}

func (this *MatchRule) SetCountryProvinceCity(country, province, city string) {
	this.Country = country
	this.Province = province
	this.City = city
}

func (this *MatchRule) SetTagId(tagId string) {
	this.TagId = tagId
}

func (this *MatchRule) SetLanguage(language string) {
	this.Language = language
}
