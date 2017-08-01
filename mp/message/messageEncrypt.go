package message

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"math/big"

	"github.com/chinasarft/wechat/mp/token"
)

var ErrAppIdNotMatch = errors.New("appidNotMatch")
var ErrEncryptLength = errors.New("EncryptLengthNotMatch")

//EncryptedXMLMsg 安全模式下的消息体
type EncryptMsg struct {
	XMLName    xml.Name `xml:"xml" json:"-"`
	ToUserName string   `xml:"ToUserName" json:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"    json:"Encrypt"`
}

//ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptMsg struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	Encrypt      CDATA    `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature CDATA    `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        CDATA    `xml:"Nonce"        json:"Nonce"`
}

func handleEncrypt(query *WechatQuery, w http.ResponseWriter,
	r *http.Request) ([]byte, error) {

	encryptMsg, err := getEncryptMsg(r)
	if err != nil {
		return nil, err
	}

	if !query.validateMsgSignature(encryptMsg.Encrypt) {
		return nil, ErrMsgSig
	}

	aesKey := token.GetAesKey()
	plainData, err := aesDecrypt(encryptMsg.Encrypt, aesKey)
	if err != nil {
		return nil, err
	}

	rawMsg, err := validateTailAppIdAndGetPlain(plainData)
	if err != nil {
		return nil, err
	}

	mixMsg, err := getMixMessage(rawMsg)
	if err != nil {
		return nil, err
	}

	result := handleMessage(mixMsg)

	respByte, err := xml.MarshalIndent(result, "", "")
	if err != nil {
		return nil, err
	}
	cipherRespByte, err := makeEncryptResponseXml(respByte)
	if err != nil {
		return nil, err
	}
	w.Header().Set("Content-Type", "text/xml")

	return cipherRespByte, nil
}

func getEncryptMsg(r *http.Request) (*EncryptMsg, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	encryptMsg := &EncryptMsg{}
	err = xml.Unmarshal(body, encryptMsg)
	if err != nil {
		return nil, err
	}
	return encryptMsg, nil
}

func validateTailAppIdAndGetPlain(plainText []byte) ([]byte, error) {

	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)

	//validate tail appid
	appId := token.GetAppId()
	appIdstart := 20 + length
	id := plainText[appIdstart : int(appIdstart)+len(appId)]
	if string(id) != appId {
		return nil, ErrAppIdNotMatch
	}

	return plainText[20 : 20+length], nil
}

func aesDecrypt(base64Data string, aesKey []byte) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, ErrEncryptLength
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

func makeEncryptResponseXml(plainData []byte) ([]byte, error) {
	// Encrypt part2: Length bytes
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(plainData)))
	if err != nil {
		fmt.Println("Binary write err:", err)
		return nil, err
	}
	bodyLength := buf.Bytes()

	// Encrypt part1: Random bytes
	randomBytes := []byte("abcdefghijklmnop")

	//with part4 - appID
	appId := token.GetAppId()

	// Msg to be encrypt Part3
	// all part need to be encrypt
	decoratePlainData := bytes.Join([][]byte{randomBytes, bodyLength, plainData, []byte(appId)}, nil)

	cipherData, err := aesEncrypt(decoratePlainData, token.GetAesKey())
	if err != nil {
		return nil, err
	}

	responseEncryptMsg := &ResponseEncryptMsg{}
	responseEncryptMsg.Encrypt = CDATA(base64.StdEncoding.EncodeToString(cipherData))

	numNonce, err := rand.Int(rand.Reader, big.NewInt(1000000000))
	if err != nil {
		return nil, err
	}
	strNonce := numNonce.String()
	responseEncryptMsg.Nonce = CDATA(strNonce)

	responseEncryptMsg.Timestamp = time.Now().Unix()
	sig := makeMsgSignature(
		strconv.Itoa(int(responseEncryptMsg.Timestamp)),
		string(responseEncryptMsg.Nonce), string(responseEncryptMsg.Encrypt))
	responseEncryptMsg.MsgSignature = CDATA(sig)

	return xml.MarshalIndent(responseEncryptMsg, "", "")
}

func aesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey)
	plainData = PKCS7Padding(plainData, k)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)

	return cipherData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
