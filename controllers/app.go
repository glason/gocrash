package controllers

import (
	"crash.android.meituan/models"
	//"fmt"
	"github.com/astaxie/beego"
	"regexp"
	"sort"
	"strconv"
	//"strings"
)

type AppController struct {
	beego.Controller
}

type CrashLog struct {
	Name        string
	Description string
	Md5         string
	Count       int
}

type CrashStats struct {
	Name  string
	Count int
}

//每页crash数目
const CRASH_PER_PAGE = 20

func (this *AppController) Get() {
	//参数
	var app, version, channel, date string
	var page, logcount, total int
	//CrashLog
	var crashLog []CrashLog
	//报表
	var dateStats, typeStats, versionStats, deviceStats []CrashStats
	//select string
	var dateOpt, versionOpt, channelOpt []string

	app = this.Ctx.Params[":app"]
	version = this.GetString("version")
	date = this.GetString("date")
	channel = this.GetString("channel")
	tpage := this.GetString("page")
	page, err := strconv.Atoi(tpage)
	if err != nil {
		page = 1
	}
	var tmpmap [][]map[string][]byte
	total, logcount, tmpmap = http.GetDataForAppPage(app, date, version, channel, (page-1)*CRASH_PER_PAGE, CRASH_PER_PAGE)

	re, _ := regexp.Compile("java:\\d+")
	for _, m := range tmpmap[0] {
		des := string(m["log"])
		count, _ := strconv.Atoi(string(m["count(*)"]))
		if len(des) > 300 {
			des = des[:300]
		}
		des = re.ReplaceAllString(des, "java:***")
		t := CrashLog{string(m["crashtype"]), des, string(m["md5"]), count}
		crashLog = append(crashLog, t)
	}

	for _, m := range tmpmap[1] {
		s := string(m["date"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		dateOpt = append(dateOpt, s)
		dateStats = append(dateStats, t)
	}
	sort.Strings(dateOpt)
	for _, m := range tmpmap[2] {
		s := string(m["crashtype"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		//append(typeOpt, s)
		typeStats = append(typeStats, t)
	}
	//sort.Strings(typeOpt)
	for _, m := range tmpmap[3] {
		s := string(m["app"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		versionOpt = append(versionOpt, s)
		versionStats = append(versionStats, t)
	}
	sort.Strings(versionOpt)
	for _, m := range tmpmap[4] {
		s := string(m["dm"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		//deviceOpt = append(deviceOpt, s)
		deviceStats = append(deviceStats, t)
	}
	//sort.Strings(deviceOpt)
	for _, m := range tmpmap[5] {
		s := string(m["ch"])
		channelOpt = append(channelOpt, s)
	}
	sort.Strings(channelOpt)

	this.TplNames = "app.html"
	this.Data["App"] = app
	this.Data["Date"] = date
	this.Data["Channel"] = channel
	this.Data["Version"] = version
	this.Data["Total"] = total
	this.Data["CrashLog"] = crashLog
	this.Data["DateStats"] = dateStats
	this.Data["VersionStats"] = versionStats
	this.Data["DeviceStats"] = deviceStats
	this.Data["TypeStats"] = typeStats
	this.Data["DateOpt"] = dateOpt
	this.Data["VersionOpt"] = versionOpt
	this.Data["ChannelOpt"] = channelOpt
	this.Data["CurPage"] = page
	this.Data["TotalPage"] = logcount/CRASH_PER_PAGE + 1
}
