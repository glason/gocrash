package controllers

import (
	"crash.android.meituan/models"
	//"fmt"
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

func (this *AppController) Get() {
	//筛选条件
	var app, version, date, channel, crashType string
	//页面form
	var allVersion, allDate, allChannel string
	//页面报表
	var mapVersion, mapDate, mapChannel, mapType map[string]int
	//crash总数
	var total int
	//解析的所有crash log
	var crashLog []CrashLog
	//log对应的crash obj
	var crashCount map[string][]http.CrashObj

	//解析参数
	tapp := this.Ctx.Params[":app"]
	tversion := this.GetString("version")
	tdate := this.GetString("date")
	tchannel := this.GetString("channel")
	tpage := this.GetString("page")
	tcrashType := this.GetString("type")
	page, err := strconv.Atoi(tpage)
	if err != nil {
		page = 1
	}

	//如果已经解析过数据，则直接使用crashLog中的数据
	app, _ = this.GetSession("app").(string)
	version, _ = this.GetSession("version").(string)
	date, _ = this.GetSession("date").(string)
	channel, _ = this.GetSession("channel").(string)
	crashType, _ = this.GetSession("crashType").(string)

	if app == tapp && version == tversion && date == tdate && channel == tchannel && crashType == tcrashType {
		allVersion, _ = this.GetSession("allVersion").(string)
		allDate, _ = this.GetSession("allDate").(string)
		allChannel, _ = this.GetSession("allChannel").(string)

		mapVersion, _ = this.GetSession("mapVersion").(map[string]int)
		mapDate, _ = this.GetSession("mapDate").(map[string]int)
		mapChannel, _ = this.GetSession("mapChannel").(map[string]int)
		mapType, _ = this.GetSession("mapType").(map[string]int)

		total, _ = this.GetSession("total").(int)
		crashLog, _ = this.GetSession("crashLog").([]CrashLog)
		crashCount = this.GetSession("crashCount").(map[string][]http.CrashObj)

	} else {
		crashObj, tmapVersion, tmapChannel, tmapDate := http.GetFilteredCrashObj(tapp, tversion, tchannel, tdate, tcrashType)
		mapVersion = tmapVersion
		mapChannel = tmapChannel
		mapDate = tmapDate
		//解析所有可选条件
		if app != tapp {
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
		} else {
			allVersion, _ = this.GetSession("allVersion").(string)
			allDate, _ = this.GetSession("allDate").(string)
			allChannel, _ = this.GetSession("allChannel").(string)
		}
		//记录下数据
		this.SetSession("app", tapp)
		this.SetSession("version", tversion)
		this.SetSession("date", tdate)
		this.SetSession("channel", tchannel)
		this.SetSession("crashType", tcrashType)

		this.SetSession("allVersion", allVersion)
		this.SetSession("allDate", allDate)
		this.SetSession("allChannel", allChannel)

		this.SetSession("mapVersion", mapVersion)
		this.SetSession("mapDate", mapDate)
		this.SetSession("mapChannel", mapChannel)

		total = len(crashObj)
		//crash log对应的crash obj
		crashCount = make(map[string][]http.CrashObj)
		this.SetSession("total", total)

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
		this.SetSession("crashCount", crashCount)
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
			mapType[name] += len(count)
			index++
		}
		this.SetSession("mapType", mapType)
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
		this.SetSession("crashLog", crashLog)
	}
	this.SetSession("page", page)

	this.TplNames = "app.html"
	this.Data["App"] = tapp
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
