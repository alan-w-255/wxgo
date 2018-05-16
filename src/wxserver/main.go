package main

import (
	_ "alanwong/wxgo/src/wxserver/routers"
	"alanwong/wxgo/src/wxserver/utils"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
}

func main() {

	// 配置日记输出文件
	logs.SetLogger("console")
	logger := logs.GetBeeLogger()

	// 更新微信的 access token
	s, err := util.GetWXAccessTokenFromWX()
	if err != nil {
		logger.Emergency("unable to get Access token: %s", err.Error())
		panic("Fatal: unbable to get Access token")
	}
	if err := beego.AppConfig.Set("wxaccesstoken", s); err != nil {
		panic("Fatal: unable to set Access token")
	}
	updateDuration, err := beego.AppConfig.Int64("wxaccesstokenupdateduration")
	if err != nil {
		panic("Fatal: unable to get wx access token update duration from configuraiton file")
	}
	timer := time.NewTicker((time.Duration)(updateDuration) * time.Second)
	util.UpdateWXAccessToken(timer) // 定时刷新 微信 access token

	defer func() {
		timer.Stop()
	}()

	beego.Run()
}
