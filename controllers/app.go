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

//每页crash数目
const CRASH_PER_PAGE = 20

//筛选条件
var app, version, date, channel string

var allVersion, allDate, allChannel string

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
		page = 1
	}
	//如果已经解析过数据，则直接使用crashLog中的数据
	if app == tapp && version == tversion && date == tdate && channel == tchannel {
		fmt.Println("use old data...")
		showPage(this, page)
		return
	}
	crashObj := crash.GetFilteredCrashObj(tapp, tversion, tchannel, tdate)
	//解析所有可选条件
	if app != tapp {
		for _, v := range crashObj {
			if strings.Index(allVersion, v.App) == -1 {
				allVersion = allVersion + "\n" + v.App
			}
			if strings.Index(allChannel, v.Ch) == -1 {
				allChannel = allChannel + "\n" + v.Ch
			}
			if strings.Index(allDate, v.Date) == -1 {
				allDate = allDate + "\n" + v.Date
			}
		}
	}
	//记录下数据
	app = tapp
	version = tversion
	date = tdate
	channel = tchannel

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
	fmt.Println(crashCount)
	fmt.Println("crashCount len:", len(crashCount))
	crashLog = make([]CrashLog, len(crashCount))
	index := 0
	for log, count := range crashCount {
		if strings.Index(log, ":") > 0 {
			crashLog[index] = CrashLog{log[:strings.Index(log, ":")], log + "...", count}
		} else {
			crashLog[index] = CrashLog{log[:50], log + "...", count}
		}
		index++
	}
	crashLog = crashLog[:index]
	for i := 0; i < index; i++ {
		max := i
		for j := i + 1; j < index; j++ {
			if crashLog[j].Count > crashLog[max].Count {
				max = j
			}
		}
		if max != i {
			tmp := crashLog[i]
			crashLog[i] = crashLog[max]
			crashLog[max] = tmp
		}
	}
	fmt.Println("crashLog len:", len(crashLog))

	showPage(this, page)
}

func showPage(this *CrashController, page int) {
	this.TplNames = "crash.html"
	this.Data["App"] = app
	this.Data["Total"] = total
	this.Data["CurPage"] = page
	this.Data["AllVersion"] = strings.Split(allVersion, "\n")[1:]
	this.Data["AllDate"] = strings.Split(allDate, "\n")[1:]
	this.Data["AllChannel"] = strings.Split(allChannel, "\n")[1:]
	var totalPage int
	totalPage = len(crashLog) / CRASH_PER_PAGE
	if len(crashLog)%CRASH_PER_PAGE != 0 {
		totalPage++
	}
	this.Data["TotalPage"] = totalPage
	this.Data["CurPage"] = page
	if (page-1)*CRASH_PER_PAGE < len(crashLog) {
		if (page-1)*CRASH_PER_PAGE+CRASH_PER_PAGE < len(crashLog) {
			this.Data["AllCrashLog"] = crashLog[(page-1)*CRASH_PER_PAGE : (page-1)*CRASH_PER_PAGE+CRASH_PER_PAGE]
		} else {
			this.Data["AllCrashLog"] = crashLog[(page-1)*CRASH_PER_PAGE : len(crashLog)]
		}
	}
}
