package main

import (
	_ "alanwong/wxgo/src/wxserver/routers"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger("console")
	logger := logs.GetBeeLogger()
	s, err := utils.GetWXAccessTokenFromWX("client_credential")
	if err != nil {
		logger.Emergency("unable to get Access token: %s", err.Error())
		panic("Fatal: unbable to get Access token")
	}
	fmt.Println("获取的access token: ", s)
	if err := beego.AppConfig.Set("wxaccesstoken", s); err != nil {
		panic("Fatal: unable to set Access token")
	}
	fmt.Println("access token:", beego.AppConfig.String("wxaccesstoken"))
	beego.Run()
}
