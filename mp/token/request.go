package token

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpPost(urlPre, contentType string, content []byte) ([]byte, error) {

	url := urlPre + GetAccessToken()

	resp, err := http.Post(url, contentType, bytes.NewReader(content))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func HttpGet(urlPre string) ([]byte, error) {

	url := urlPre + GetAccessToken()

	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
