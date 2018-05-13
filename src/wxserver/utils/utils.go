package util

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
func GetWXAccessTokenFromWX() (string, error) {
	tpl := beego.AppConfig.String("wxapiaccesstoken")
	appid := beego.AppConfig.String("appid")
	appsecret := beego.AppConfig.String("appsecret")
	grantType := "client_credential"
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
		l.Info("get %s response status code: 200", urlstr)
		var accessTokRes AccessTokenResponse
		var errRes ErrorResponse
		if err := json.Unmarshal(body, &accessTokRes); err == nil {
			return accessTokRes.AccessToken, nil
		} else if err := json.Unmarshal(body, &errRes); err == nil {
			return "", fmt.Errorf("get %s response errcode: %d, errmsg: %s", urlstr, errRes.Errcode, errRes.Errmsg)
		} else {
			return "", fmt.Errorf("unable to unmashal %s response: %s error: %s", urlstr, string(body), err.Error())
		}
	}
	return "", fmt.Errorf("error response status code: %d", res.StatusCode)

}

//GetWXAccessTokenFromLocal 从本地缓存获取access token
func GetWXAccessTokenFromLocal() string {
	return beego.AppConfig.String("wxaccesstoken")
}

// UpdateWXAccessToken 定时刷新微信access token
func UpdateWXAccessToken(timer *time.Ticker) {
	l := logs.GetBeeLogger()
	go func() {
		// fmt.Println("start tick to update access token")
		for range timer.C {
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second)
				s, err := GetWXAccessTokenFromWX()
				// fmt.Printf("update access_tokenn: %s\n", s)
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

// CreateMenu 创建自定义菜单
func CreateMenu() error {
	urlTpl := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
	accessToke := GetWXAccessTokenFromLocal()

	url := fmt.Sprintf(urlTpl, accessToke)
	menuconfigpath := beego.AppConfig.String("wxmenuconfigpath")
	config, err := os.Open(menuconfigpath)
	logger := logs.GetBeeLogger()
	if err != nil {
		logger.Error("error occured while open %s", menuconfigpath)
		logger.Error(err.Error())
		errmsg := fmt.Sprintf("打开文件 %s 出错", menuconfigpath)
		return errors.New(errmsg)
	}
	w, _ := ioutil.ReadAll(config)
	res, err := http.Post(url, "application/json", strings.NewReader(string(w)))
	if err != nil {
		logger.Error("post 请求 %s 时出错!", url)
		logger.Error(err.Error())
		return errors.New("post 请求出错")
	}
	wxres := struct {
		ErrorCode int    `json:"errcode"`
		Errmsg    string `json:"errmsg"`
	}{}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		logger.Error("读取post %s 请求返回的数据出错", url)
		return errors.New("读取请求返回数据出错")
	}
	if err := json.Unmarshal(data, &wxres); err != nil {
		errmsg := fmt.Sprintf("解析post %s 返回数据出错", url)
		logger.Info("statu code : %d", res.StatusCode)
		logger.Error(err.Error())
		logger.Error(errmsg)
		logger.Error("返回数据为: %s", string(data))
		return errors.New(errmsg)
	}
	if wxres.ErrorCode == 0 {
		logger.Info("创建菜单成功")
		return nil
	} else {
		logger.Error("创建菜单出错, errcode: %d, errmsg: %s", wxres.ErrorCode, wxres.Errmsg)
		return errors.New("创建菜单出错")
	}
}
