package controllers

import (
	"github.com/astaxie/beego"
	"net/url"
	"strings"
)

type CrashController struct {
	beego.Controller
}

func (this *CrashController) Get() {
	this.TplNames = "crash.html"
	this.Data["Appnm"] = App

	crash, _ := url.QueryUnescape(this.Ctx.Params[":crash"])

	crashObj := CrashCount[crash]
	this.Data["Total"] = len(crashObj)
	this.Data["CrashObj"] = crashObj
	if len(crashObj) >= 1 {
		log := crashObj[0].Log
		if strings.Index(log, ":") > 0 {
			this.Data["CrashType"] = log[:strings.Index(log, ":")]
		}
		this.Data["CrashDetail"] = log
	}
	//报表
	mapVersion := make(map[string]int)
	mapDate := make(map[string]int)
	mapChannel := make(map[string]int)
	for _, v := range crashObj {
		mapVersion[v.App] += 1
		mapChannel[v.Ch] += 1
		mapDate[v.Date] += 1
	}
	this.Data["MapVersion"] = mapVersion
	this.Data["MapDate"] = mapDate
	this.Data["MapChannel"] = mapChannel
}
