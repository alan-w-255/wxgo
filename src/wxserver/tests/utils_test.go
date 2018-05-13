package test

import (
	util "alanwong/wxgo/src/wxserver/utils"
	"fmt"
	"testing"
	"time"

	"github.com/astaxie/beego"
)

func TestGetWXAccessTokenFromWX(t *testing.T) {
	s, _ := util.GetWXAccessTokenFromWX()
	if len(s) == 0 {
		t.FailNow()
	}
}

func TestGetWXAccessTokenFromLocal(t *testing.T) {
	updateDuration, err := beego.AppConfig.Int64("wxaccesstokenupdateduration")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	timer := time.NewTicker((time.Duration)(updateDuration))
	util.UpdateWXAccessToken(timer)
	time.Sleep(time.Second * 3)
	s := util.GetWXAccessTokenFromLocal()
	if len(s) == 0 {
		t.Error("get an empty access token")
		t.FailNow()
	}
	timer.Stop()
}

func TestUpdateWXAccessToken(t *testing.T) {
	updateDuration, err := beego.AppConfig.Int64("wxaccesstokenupdateduration")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	timer := time.NewTicker((time.Duration)(updateDuration))
	util.UpdateWXAccessToken(timer)
	time.Sleep(time.Second * 3)
	s := util.GetWXAccessTokenFromLocal()
	if len(s) == 0 {
		t.Error("get an empty access token string")
		t.FailNow()
	}
	timer.Stop()
}

func TestCreateMenu(t *testing.T) {
	accessToken, err := util.GetWXAccessTokenFromWX()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	beego.AppConfig.Set("wxaccesstoken", accessToken)
	if err := util.CreateMenu(); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}
