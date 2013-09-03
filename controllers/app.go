package controllers

import (
	"crash.android.meituan/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type CrashController struct {
	beego.Controller
}

type CrashLog struct {
	Name        string
	Description string
	Count       int
}

//筛选条件
var app, version, date, channel string

//crash总数
var total int

//解析的所有crash log
var crashLog []CrashLog

func (this *CrashController) Get() {

	tapp := this.Ctx.Params[":app"]
	tversion := this.GetString("version")
	tdate := this.GetString("date")
	tchannel := this.GetString("channel")
	tpage := this.GetString("page")
	//分页展示，每页10条
	page, err := strconv.Atoi(tpage)
	if err != nil {
		page = 0
	}
	//如果已经解析过数据，则直接使用crashLog中的数据
	if app == tapp && version == tversion && date == tdate && channel == tchannel {
		fmt.Println("use old data...")
		showPage(this, page)
		return
	}
	//记录下数据
	app = tapp
	version = tversion
	date = tdate
	channel = tchannel

	crashObj := crash.GetFilteredCrashObj(app, version, channel, date)
	fmt.Println("crashObj len:", len(crashObj))
	total = len(crashObj)
	crashCount := make(map[string]int)

	for _, v := range crashObj {
		var log string
		//取Log前300字符作为description
		if len(v.Log) >= 300 {
			log = v.Log[:300]
		} else {
			log = v.Log
		}
		crashCount[log]++
	}
	crashLog = make([]CrashLog, len(crashCount))
	index := 0
	for log, count := range crashCount {
		maxK := log
		maxV := count
		for k, v := range crashCount {
			if v > maxV {
				maxK = k
				maxV = v
			}
		}
		if strings.Index(maxK, ":") > 0 {
			crashLog[index] = CrashLog{maxK[:strings.Index(maxK, ":")], maxK + "...", maxV}
		} else {
			crashLog[index] = CrashLog{maxK[:50], maxK + "...", maxV}
		}
		index++
		delete(crashCount, maxK)
	}

	showPage(this, page)
}

func showPage(this *CrashController, page int) {
	this.TplNames = "crash.html"
	this.Data["App"] = app
	this.Data["Total"] = total
	fmt.Println("show page:", page)
	if page*10 < len(crashLog) {
		if page*10+10 < len(crashLog) {
			this.Data["AllCrashLog"] = crashLog[page*10 : page*10+10]
			fmt.Println("crashLog[", page*10, ":", page*10+10, "]")
		} else {
			this.Data["AllCrashLog"] = crashLog[page*10 : len(crashLog)]
		}
	}
}
