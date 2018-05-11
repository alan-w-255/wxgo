package routers

import (
	"alanwong/wxgo/src/wxserver/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/wx", &controllers.WXController{})
}
