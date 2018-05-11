package controllers

import (
	util "alanwong/wxgo/src/wxserver/utils"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// MainController /
type MainController struct {
	beego.Controller
}

// Get 处理 get("/")
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

// WXController 处理"/wx"路由下 逻辑
type WXController struct {
	beego.Controller
}

// Get 处理get("/wx")逻辑
func (c *WXController) Get() {
	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")
	token := beego.AppConfig.String("wxtoken")

	l := logs.GetBeeLogger()
	if util.IsWXServer(signature, timestamp, nonce, echostr, token) {
		l.Info("get request from wx server")
		c.Ctx.WriteString(echostr)
	} else {

		l.Debug("get request from unknow server")
	}
}
