package controllers

import (
	"crash.android.meituan/models"
	"github.com/astaxie/beego"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type AppController struct {
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
var App, version, date, channel, crashType string

var allVersion, allDate, allChannel string
var mapVersion, mapDate, mapChannel, mapType map[string]int

//crash总数
var total int
var cur_page int

//解析的所有crash log
var crashLog []CrashLog

//log对应的crash obj
var crashCount map[string][]http.CrashObj

func (this *AppController) Get() {

	tapp := this.Ctx.Params[":app"]
	tversion := this.GetString("version")
	tdate := this.GetString("date")
	tchannel := this.GetString("channel")
	tpage := this.GetString("page")
	tcrashType := this.GetString("type")
	//分页展示，每页10条
	page, err := strconv.Atoi(tpage)
	if err != nil {
		page = 1
	}
	//如果已经解析过数据，则直接使用crashLog中的数据
	if App == tapp && version == tversion && date == tdate && channel == tchannel && crashType == tcrashType {
		showPage(this, page)
		return
	}
	crashObj, tmapVersion, tmapChannel, tmapDate := http.GetFilteredCrashObj(tapp, tversion, tchannel, tdate, tcrashType)
	mapVersion = tmapVersion
	mapChannel = tmapChannel
	mapDate = tmapDate
	//解析所有可选条件
	if App != tapp {
		allVersion = ""
		allChannel = ""
		allDate = ""
		for k, _ := range mapVersion {
			allVersion = allVersion + "\n" + k
		}
		for k, _ := range mapChannel {
			allChannel = allChannel + "\n" + k
		}
		for k, _ := range mapDate {
			allDate = allDate + "\n" + k
		}
	}
	//记录下数据
	App = tapp
	version = tversion
	date = tdate
	channel = tchannel
	crashType = tcrashType

	total = len(crashObj)
	crashCount = make(map[string][]http.CrashObj)

	for _, v := range crashObj {
		var log string
		//取Log前300字符作为description
		if len(v.Log) >= 300 {
			log = v.Log[:300]
		} else {
			log = v.Log
		}
		crashCount[log] = append(crashCount[log], v)
	}
	crashLog = make([]CrashLog, len(crashCount))
	index := 0
	mapType = make(map[string]int)
	re, _ := regexp.Compile("^\\w+\\.\\w+\\.\\w+")
	for log, count := range crashCount {
		name := re.FindString(log)
		if name == "" {
			name = log[:50]
		}
		crashLog[index] = CrashLog{name, log, len(count)}
		mapType[name] += 1
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

	showPage(this, page)
}

func getCrashByIndex(index int) []http.CrashObj {
	i := index + (cur_page-1)*CRASH_PER_PAGE
	if i < len(crashLog) {
		log := crashLog[i].Description
		return crashCount[log]
	} else {
		return nil
	}
}

func showPage(this *AppController, page int) {
	cur_page = page
	this.TplNames = "app.html"
	this.Data["App"] = App
	this.Data["Total"] = total

	tmp := strings.Split(allVersion, "\n")[1:]
	sort.Strings(tmp)
	this.Data["AllVersion"] = tmp
	this.Data["MapVersion"] = mapVersion

	tmp = strings.Split(allDate, "\n")[1:]
	sort.Strings(tmp)
	this.Data["AllDate"] = tmp
	this.Data["MapDate"] = mapDate

	tmp = strings.Split(allChannel, "\n")[1:]
	sort.Strings(tmp)
	this.Data["AllChannel"] = tmp
	this.Data["MapChannel"] = mapChannel

	this.Data["MapType"] = mapType

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
