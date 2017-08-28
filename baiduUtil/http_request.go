package baiduUtil

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// HTTPGet 简单实现 http 访问 GET 请求
func HTTPGet(urlStr string) (body []byte, err error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//Fetch 实现 http／https 访问 和 GET／POST 请求，根据给定的 urlStr (网址), jar (cookie容器 指针), post (post数据 指针), header (请求头数据 指针) 进行网站访问。
//返回值分别为 网站主体, 错误
func Fetch(urlStr string, jar *cookiejar.Jar, post, header *map[string]string) (body []byte, err error) {
	URL, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalln(err)
	}

	var req *http.Request
	httpClient := &http.Client{Timeout: 3e10} // 30s
	if jar != nil {
		httpClient.Jar = jar
	}

	if URL.Scheme == "https" {
		httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	if post == nil {
		req, _ = http.NewRequest("GET", urlStr, nil)
		addHeader(req, header)
	} else {
		query := url.Values{}
		for k, v := range *post {
			query.Set(k, v)
		}
		req, _ = http.NewRequest("POST", urlStr, strings.NewReader(query.Encode()))
		addHeader(req, header)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func addHeader(req *http.Request, header *map[string]string) {
	if header != nil {
		for Header, HeaderMessage := range *header {
			req.Header.Add(Header, HeaderMessage)
		}
	}
}