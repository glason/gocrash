package controllers

import (
	"crash.android.meituan/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.TplNames = "index.html"
	this.Data["Time"] = http.UpdateTime
}
