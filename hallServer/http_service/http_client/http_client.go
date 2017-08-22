package http_client

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/lovelly/leaf/log"
)

func NewHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}

	return &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
		//Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
	}
}

func PostJSON(url string, data interface{}) *simplejson.Json {
	log.Debug("post:url :%s, data:%v", url, data)
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		log.Error("at PostJSON  json Marshal err :%s", err.Error())
		return nil
	}

	request, _ := http.NewRequest("POST", url, bytes.NewReader(dataBytes))
	request.Header.Add("Content-Type", "application/json")
	httpClient := NewHttpClient()
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Error("at PostJSON resp error err :%s,  resp :%s:", err.Error(), resp)
		return nil
	}
	return getJsonResponse(resp)
}

func PostForm(url string, data url.Values) *simplejson.Json {
	log.Debug("post:url :%s, data:%v", url, data)
	httpClient := NewHttpClient()
	resp, err := httpClient.PostForm(url, data)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return getJsonResponse(resp)
}

func GetUrl(url string, data url.Values) *simplejson.Json {
	log.Debug("%s, data:%v", url, data)
	httpClient := NewHttpClient()
	resp, err := httpClient.Get(fmt.Sprintf("%s?%s", url, data.Encode()))
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return getJsonResponse(resp)
}

func getResponse(resp *http.Response) []byte {
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	log.Debug("getResponse:%v", resp)
	var result []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ := gzip.NewReader(resp.Body)
		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)

			if err != nil && err != io.EOF {
				panic(err)
			}

			if n == 0 {
				break
			}
			result = append(result, buf...)
		}
	default:
		result, _ = ioutil.ReadAll(resp.Body)
	}

	log.Debug("http respCode:%v, result:%v", resp.StatusCode, string(result))
	return result
}

func getJsonResponse(resp *http.Response) *simplejson.Json {
	result := getResponse(resp)
	jsonResponse, err := simplejson.NewJson(result)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return jsonResponse
}
