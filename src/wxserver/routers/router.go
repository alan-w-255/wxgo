package routers

import (
	"alanwong/wxgo/src/wxserver/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/wx", &controllers.WXController{})

	// todo: 开发绑定用户功能
	beego.Router("/wx/usrbind", &controllers.BindUsrController{})
}
