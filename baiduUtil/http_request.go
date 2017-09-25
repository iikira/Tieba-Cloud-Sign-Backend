package baiduUtil

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	//IsGzip 是否启用Gzip
	IsGzip = true
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
func Fetch(method string, urlStr string, jar *cookiejar.Jar, post interface{}, header map[string]string) (body []byte, err error) {
	var (
		req   *http.Request
		obody io.Reader
	)
	httpClient := &http.Client{Timeout: 3e10} // 30s
	if jar != nil {
		httpClient.Jar = jar
	}

	if HTTPSRE.MatchString(urlStr) {
		httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	if _, ok := header["Accept-Encoding"]; !ok && header != nil && IsGzip {
		header["Accept-Encoding"] = "gzip"
	}

	if post != nil {
		switch postValue := post.(type) {
		case map[string]string:
			query := url.Values{}
			for k, v := range postValue {
				query.Set(k, v)
			}
			obody = strings.NewReader(query.Encode())
		case string:
			obody = strings.NewReader(postValue)
		}
	}
	req, err = http.NewRequest(method, urlStr, obody)
	if err != nil {
		return nil, err
	}
	addHeader(req, header)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ = ioutil.ReadAll(resp.Body)
	if IsGzip {
		undatas, err := DecompressGZIP(bytes.NewReader(body))
		if err == nil {
			return undatas, nil
		}
	}
	resp.Body.Close()
	return body, nil
}

func addHeader(req *http.Request, header map[string]string) {
	if header != nil {
		for Header, HeaderMessage := range header {
			req.Header.Add(Header, HeaderMessage)
		}
	}
}
