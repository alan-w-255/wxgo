package util

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// AccessTokenResponse 微信返回的access token
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`

	ExpiresIn int `json:"expires_in"`
}

// ErrorResponse 微信返回的错误码
type ErrorResponse struct {
	Errcode int
	Errmsg  string
}

// IsWXServer 判断请求是不是来自微信后台
func IsWXServer(signature string, timestamp string, nonce string, echostr string, wxtoken string) bool {

	tmpSlice := []string{wxtoken, timestamp, nonce}
	sort.Strings(tmpSlice)
	tmpStr := strings.Join(tmpSlice, "")

	hasher := sha1.New()
	hasher.Write(([]byte)(tmpStr))
	sha1str := hex.EncodeToString(hasher.Sum(nil))

	if sha1str == signature {
		return true
	}
	return false

}

// GetWXAccessTokenFromWX 从微信后台获取微信access token
func GetWXAccessTokenFromWX(grantType string) (string, error) {
	tpl := beego.AppConfig.String("wxapiaccesstoken")
	appid := beego.AppConfig.String("appid")
	appsecret := beego.AppConfig.String("appsecret")
	urlstr := fmt.Sprintf(tpl, grantType, appid, appsecret)

	l := logs.GetBeeLogger()

	res, err := http.Get(urlstr)
	if err != nil {
		l.Error("请求获取access token 出错: ")
		l.Error(err.Error())
		return "", nil
	}

	body := make([]byte, 1024)
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		l.Error("请求获取access token 出错: ")
		l.Error(err.Error())
	}
	if res.StatusCode == 200 {
		l.Info("response status code: 200")
		var accessTokRes AccessTokenResponse
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &accessTokRes); err == nil {
			return accessTokRes.AccessToken, nil
		} else if err := json.Unmarshal(body, &errRes); err == nil {
			return "", fmt.Errorf("errcode: %d, errmsg: %s", errRes.Errcode, errRes.Errmsg)
		} else {
			return "", fmt.Errorf("error: %s", err.Error())
		}
	}
	return "", fmt.Errorf("error response status code: %d", res.StatusCode)

}

//GetWXAccessTokenFromLocal 从本地缓存获取access token
func GetWXAccessTokenFromLocal(grantType string) string {
	return beego.AppConfig.String("wxaccesstoken")
}

// UpdateWXAccessToken 定时刷新微信access token
func UpdateWXAccessToken(grantType string) {
	updateduration, _ := beego.AppConfig.Int64("wxaccesstokenupdateduration")
	timer := time.NewTicker((time.Duration)(updateduration) * time.Second)
	l := logs.GetBeeLogger()
	go func() {
		for range timer.C {
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second)
				s, err := GetWXAccessTokenFromWX("client_credential")
				if err != nil {
					l.Error(err.Error())
					if i == 3 {
						timer.Stop()
						panic("unable to get access token")
					}
				} else {
					if err := beego.AppConfig.Set("wxaccesstoken", s); err != nil {
						l.Emergency("unable to set access token")
					}
					l.Info("update access token successfully")
					break
				}
			}
		}
	}()
}
