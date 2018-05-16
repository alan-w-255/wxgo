package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// BindUsrController 绑定用户的控制器
type BindUsrController struct {
	beego.Controller
}

// Get 请求绑定用户网页
func (c *BindUsrController) Get() {
	c.Data["website"] = "用户绑定"
	c.TplName = "bindusr.html"

}

// Post 验证用户提交绑定用户表单, 绑定用户
func (c *BindUsrController) Post() {
	logger := logs.GetBeeLogger()

	fmt.Printf("获取到的json 数据: %s\n", c.Ctx.Input.RequestBody)

	type user struct {
		Name      interface{} `form:"user_name"`
		IDCardNum interface{} `form:"id_card_number"`
		PhoneNum  string      `form:"tel_number"`
	}
	u := user{}

	if err := c.ParseForm(&u); err != nil {
		logger.Error("解析表单出错")
		logger.Error(err.Error())
		logger.Debug("获取到的表单数据: %s", c.GetString("user_name"))
	}

	c.Data["website"] = "用户绑定"
	c.TplName = "bindusr.html"
}
