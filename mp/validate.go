package mp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

func ValidateWechatServer(token, timestamp, nonce, signature, echostr string) bool {
	s := []string{token, timestamp, nonce}
	sort.Sort(sort.StringSlice(s)) //将token、timestamp、nonce三个参数进行字典序排序
	s0 := strings.Join(s, "")      //将三个参数字符串拼接成一个字符串
	t := sha1.New()                //sha1加密
	io.WriteString(t, s0)
	s1 := fmt.Sprintf("%x", t.Sum(nil))
	if signature == s1 { //与signature对比
		return true
	}
	return false
}
