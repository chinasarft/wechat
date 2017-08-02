package menu

import (
	"fmt"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/chinasarft/wechat/mp/message"
	"github.com/chinasarft/wechat/mp/token"
)

const (
	MAX_FIRST_CLASS_MENU  = 3
	MAX_SECOND_CLASS_MENU = 5

	ButtonTypeView        = "view"
	ButtonTypeClick       = "click"
	ButtonTypeMiniprogram = "miniprogram"

	// 下面的按钮类型仅支持微信 iPhone5.4.1 以上版本, 和 Android5.4 以上版本的微信用户,
	// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
	ButtonTypeScancodePush    = "scancode_push"      // 扫码推事件
	ButtonTypeScancodeWaitMsg = "scancode_waitmsg"   // 扫码带提示
	ButtonTypePicSysphoto     = "pic_sysphoto"       // 系统拍照发图
	ButtonTypePicPhotoOrAlbum = "pic_photo_or_album" // 拍照或者相册发图
	ButtonTypePicWeixin       = "pic_weixin"         // 微信相册发图
	ButtonTypeLocationSelect  = "location_select"    // 发送位置

)

var ErrMaxMenu = errors.New("max menu")

type Menu struct {
	Buttons []*Button `json:"button,omitempty"`
}

type Button struct {
	Type       string    `json:"type,omitempty"`     // 非必须; 菜单的响应动作类型
	Name       string    `json:"name,omitempty"`     // 必须;  菜单标题
	Key        string    `json:"key,omitempty"`      // 非必须; 菜单KEY值, 用于消息接口推送
	Url        string    `json:"url,omitempty"`      // 非必须; 网页链接, 用户点击菜单可打开链接
	MediaId    string    `json:"media_id,omitempty"` // 非必须; 调用新增永久素材接口返回的合法media_id
	AppId      string    `json:"appid,omitempty"`
	PagePath   string    `json:"pagepath,omitempty"`
	SubButtons []*Button `json:"sub_button,omitempty"` // 非必须; 二级菜单数组
}

type MenuButton Button

func NewMenu() *Menu {
	return &Menu{}
}

func (this *Menu) AddButton(button *Button) error {
	if len(this.Buttons) >= MAX_FIRST_CLASS_MENU {
		return ErrMaxMenu
	}
	this.Buttons = append(this.Buttons, button)
	return nil
}

func (this *Menu) AddMenuButton(menuButton *MenuButton) error {
	if len(this.Buttons) >= MAX_FIRST_CLASS_MENU {
		return ErrMaxMenu
	}
	button := (*Button)(menuButton)
	this.Buttons = append(this.Buttons, button)
	return nil
}

func (this *Menu) GetJsonByte() ([]byte, error) {
	return json.MarshalIndent(this, "", "\t")
}

func (this *Menu) Submit() error {

	text, err := this.GetJsonByte()
	if err != nil {
		return err
	}

	accessToken := token.GetAccessToken()
	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/menu/create?access_token="+accessToken,
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

	result := &message.RequstWechatResult{}
	if err = json.Unmarshal(body, result); err != nil {
		return err
	}
	if result.Errcode != 0 {
		// TODO just return ErrMsg?
		return fmt.Errorf(string(body))
	}
	return nil

}

func NewClickButton(name, key string) *Button {
	return &Button{Type: ButtonTypeClick, Name: name, Key: key}
}

func NewViewButton(name, url string) *Button {
	return &Button{Type: ButtonTypeView, Name: name, Url: url}
}

func NewProgramButton(name, url, appid, pagepath string) *Button {
	return &Button{
		Type:     ButtonTypeMiniprogram,
		Name:     name,
		Url:      url,
		AppId:    appid,
		PagePath: pagepath,
	}
}

func NewScancodePushButton(name, key string) *Button {
	return &Button{Type: ButtonTypeScancodePush, Name: name, Key: key}
}

func NewScancodeWaitmsgButton(name, key string) *Button {
	return &Button{Type: ButtonTypeScancodePush, Name: name, Key: key}
}

func NewPicSysphotoButton(name, key string) *Button {
	return &Button{Type: ButtonTypePicSysphoto, Name: name, Key: key}
}

func NewPicPhotoOrAlbumButton(name, key string) *Button {
	return &Button{Type: ButtonTypePicPhotoOrAlbum, Name: name, Key: key}
}

func NewPicWeixinButton(name, key string) *Button {
	return &Button{Type: ButtonTypePicPhotoOrAlbum, Name: name, Key: key}
}

func NewLocationSelectButton(name, key string) *Button {
	return &Button{Type: ButtonTypeLocationSelect, Name: name, Key: key}
}

func NewMenuButton(name string) *MenuButton {
	return &MenuButton{Name: name}
}

func (this *MenuButton) AddSubButton(button *Button) error {
	if len(this.SubButtons) >= MAX_SECOND_CLASS_MENU {
		return ErrMaxMenu
	}
	this.SubButtons = append(this.SubButtons, button)
	return nil
}
