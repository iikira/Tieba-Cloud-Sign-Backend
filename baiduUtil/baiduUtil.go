// Package baiduUtil 为操作百度贴吧签到等功能而写的一些简易工具
package baiduUtil

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

/*
BeijingTimeOption 根据给定的 get 返回时间格式.

	get:        时间格式

	"Refer":    2017-7-21 12:02:32.000
	"printLog": 2017-7-21_12:02:32
	"day":      21
	"ymd":      2017-7-21
	"hour":     12
	默认时间戳:   1500609752
*/
func BeijingTimeOption(get string) string {
	//获取北京（东八区）时间
	CSTLoc := time.FixedZone("CST", 8*3600) // 东8区
	now := time.Now().In(CSTLoc)
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	millisecond := now.Nanosecond() / 1e6
	switch get {
	case "Refer":
		return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d.%03d", year, mon, day, hour, min, sec, millisecond)
	case "printLog":
		return fmt.Sprintf("%d-%d-%d_%02dh%02dm%02ds", year, mon, day, hour, min, sec)
	case "day":
		return fmt.Sprintf("%d", day)
	case "ymd":
		return fmt.Sprintf("%d-%d-%d", year, mon, day)
	case "hour":
		return fmt.Sprintf("%d", hour)
	default:
		return fmt.Sprintf("%d", time.Now().Unix())
	}
}

// GetURLCookieString 返回cookie字串
func GetURLCookieString(urlString string, jar *cookiejar.Jar) string {
	url, _ := url.Parse(urlString)
	cookies := jar.Cookies(url)
	cookieString := ""
	for _, v := range cookies {
		cookieString += v.String() + "; "
	}
	cookieString = strings.TrimRight(cookieString, "; ")
	return cookieString
}

// TiebaClientSignature 根据给定贴吧客户端的 post (post数据指针) 进行签名, 以通过百度服务器验证。返回值为: sign 签名字符串
func TiebaClientSignature(post *map[string]string) {
	if *post == nil {
		return
	}
	// 预设
	(*post)["_client_type"] = "2"
	(*post)["_client_version"] = "6.9.2.1"
	(*post)["_phone_imei"] = "860983036542682"
	(*post)["from"] = "mini_ad_wandoujia"
	(*post)["model"] = "HUAWEI NXT-AL10"
	(*post)["cuid"] = "61464018582906C485355A89D105ECFB|286245630389068"
	var keys []string
	for key := range *post {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	var bb bytes.Buffer
	for _, key := range keys {
		bb.WriteString(key + "=" + (*post)[key])
	}
	bb.WriteString("tiebaclient!!!")
	(*post)["sign"] = Md5Encrypt(bb.Bytes())
}

// TiebaClientRawQuerySignature 给 rawQuery 进行贴吧客户端签名
func TiebaClientRawQuerySignature(rawQuery string) string {
	return rawQuery + "&sign=" + Md5Encrypt(strings.Replace(rawQuery, "&", "", -1)+"tiebaclient!!!")
}

// Md5Encrypt 对 str 进行md5加密, 返回值为 str 加密后的密文
func Md5Encrypt(str interface{}) string {
	md5Ctx := md5.New()
	switch value := str.(type) {
	case string:
		md5Ctx.Write([]byte(str.(string)))
	case *string:
		md5Ctx.Write([]byte(*str.(*string)))
	case []byte:
		md5Ctx.Write(str.([]byte))
	case *[]byte:
		md5Ctx.Write(*str.(*[]byte))
	default:
		fmt.Println("MD5Encrypt: unknown type:", value)
		return ""
	}
	return fmt.Sprintf("%X", md5Ctx.Sum(nil))
}

// PrintErrIfExist 简易错误处理, 如果 err 存在, 就只向屏幕输出 err 。
func PrintErrIfExist(err error) {
	if err != nil {
		log.Println(err)
	}
}

// PrintErrAndExit 简易错误处理, 如果 err 存在, 向屏幕输出 err 并退出, annotate 是加在 err 之前的注释信息。
func PrintErrAndExit(annotate string, err error) {
	if err != nil {
		log.Println(annotate, err)
		os.Exit(1)
	}
}
