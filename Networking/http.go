package Networking

import (
	"io"
	"net/http"
	"net/url"
)

type Http struct {
}

func NewHttp() *Http {
	return &Http{}
}

func (h *Http) Get(urlStr string, paramsMap map[string]string) string {
	params := url.Values{}
	reqUrl, err := url.Parse(urlStr) //请求地址
	if err != nil {
		panic(err.Error())
	}
	for key, value := range paramsMap {
		params.Set(key, value) //设置参数
	}
	reqUrl.RawQuery = params.Encode()      //组合url
	resp, err := http.Get(reqUrl.String()) //发起get请求
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()            //关闭请求
	body, err := io.ReadAll(resp.Body) //解析请求信息
	if err != nil {
		panic(err.Error())
	}
	return string(body)
}
