package message

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/chinasarft/wechat/mp/token"
)

var ErrSig = errors.New("invalid signature")
var ErrMsgSig = errors.New("invalid msgSignature")

type WechatQuery struct {
	Signature    string
	Timestamp    string
	Nonce        string
	OpenId       string
	EncryptType  string
	MsgSignature string
	Echostr      string
	IsSafeMode   bool
	SignatureOk  bool
}

func getWechatQuery(r *http.Request) *WechatQuery {

	r.ParseForm()

	query := &WechatQuery{}
	query.Timestamp = strings.Join(r.Form["timestamp"], "")
	query.Nonce = strings.Join(r.Form["nonce"], "")
	query.Signature = strings.Join(r.Form["signature"], "")
	query.EncryptType = strings.Join(r.Form["encrypt_type"], "")
	query.MsgSignature = strings.Join(r.Form["msg_signature"], "")
	query.OpenId = strings.Join(r.Form["openid"], "")
	query.Echostr = strings.Join(r.Form["echostr"], "")

	if query.EncryptType == "aes" {
		query.IsSafeMode = true
	}

	return query
}

func (this *WechatQuery) validateSignature() bool {

	sig := makeSignature(this.Timestamp, this.Nonce)
	if this.Signature == sig {
		this.SignatureOk = true
	} else {
		this.SignatureOk = false
	}
	return this.SignatureOk
}

func (this *WechatQuery) validateMsgSignature(msgEncrypt string) bool {

	msgSig := makeMsgSignature(this.Timestamp, this.Nonce, msgEncrypt)
	if this.MsgSignature == msgSig {
		return true
	}
	return false
}

func makeMsgSignature(timestamp, nonce, msgEncrypt string) string {
	validateToken := token.GetValidateToken()
	s := []string{validateToken, timestamp, nonce, msgEncrypt}
	sort.Strings(s)
	s0 := strings.Join(s, "")
	t := sha1.New()
	io.WriteString(t, s0)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func makeSignature(timestamp, nonce string) string {
	validateToken := token.GetValidateToken()
	s := []string{validateToken, timestamp, nonce}
	//token timestamp nonce make dic sort
	sort.Sort(sort.StringSlice(s))
	s0 := strings.Join(s, "")
	t := sha1.New()
	io.WriteString(t, s0)
	return fmt.Sprintf("%x", t.Sum(nil))
}
