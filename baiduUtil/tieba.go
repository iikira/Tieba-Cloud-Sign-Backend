package baiduUtil

import (
	"bytes"
	"github.com/bitly/go-simplejson"
	"regexp"
	"strings"
)

// GetTiebaFid 获取贴吧fid值
func GetTiebaFid(tiebaName string) (fid string, err error) {
	b, err := HTTPGet("http://tieba.baidu.com/f/commit/share/fnameShareApi?ie=utf-8&fname=" + tiebaName)
	if err != nil {
		return
	}
	c := regexp.MustCompile(`"data":{"fid":(.*?),"can_send_pics":`).FindSubmatch(b)
	if len(c) > 1 {
		fid = string(c[1])
	}
	return
}

// GetBaiduUID 获取百度 user_id
func GetBaiduUID(name string) (uid string, err error) {
	b, err := HTTPGet("http://tieba.baidu.com/home/get/panel?un=" + name)
	if err != nil {
		return
	}
	json, err := simplejson.NewJson(b)
	if err != nil {
		return
	}
	dataJSON := json.Get("data")
	bytesUID, err := dataJSON.Get("id").MarshalJSON()
	if err != nil {
		return
	}
	uid = string(bytesUID)
	return
}

// GetBaiduNameShow 获取百度昵称
func GetBaiduNameShow(uid string) (nameShow string, err error) {
	rawQuery := "has_plist=0&need_post_count=1&rn=1&uid=" + uid
	sign := Md5Encrypt(strings.Replace(rawQuery, "&", "", -1) + "tiebaclient!!!")
	urlStr := "http://c.tieba.baidu.com/c/u/user/profile?" + rawQuery + "&sign=" + sign
	b, err := HTTPGet(urlStr)
	if err != nil {
		return
	}
	json, err := simplejson.NewJson(b)
	if err != nil {
		return
	}
	nameShow = json.GetPath("user", "name_show").MustString()
	return
}

// IsTiebaExist 检测贴吧是否存在
func IsTiebaExist(tiebaName string) bool {
	b, err := HTTPGet("http://c.tieba.baidu.com/mo/q/m?tn4=bdKSW&sub4=&word=" + tiebaName)
	PrintErrIfExist(err)
	return !bytes.Contains(b, []byte(`class="tip_text2">欢迎创建此吧，和朋友们在这里交流</p>`))
}
