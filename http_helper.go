package common

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)



func Post(url string, params map[string]string) (string, error) {

	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range params {
		err := writer.WriteField(k, v)
		if err != nil {
			return "", err
		}
	}
	err := writer.Close()
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// todo: 返回 logic 的错误
	return string(body), nil
}



