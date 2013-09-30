package controllers

import (
	"crash.android.meituan/models"
	//"fmt"
	"github.com/astaxie/beego"
	//"regexp"
	//"sort"
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

//每页crash数目
const CRASH_PER_PAGE = 20

func (this *AppController) Get() {
	//参数
	var app, version, channel, date string
	var page int
	//all crash
	var dbcrash []http.Dbcrash
	//CrashLog
	var crashLog []CrashLog
	//报表
	var mapDate, mapType, mapVersion, mapChannel, mapDevice, mapOs map[string]int
	//select string
	var allDate, allType, allVersion, allChannel, allDevice, allOs []string

	app = this.Ctx.Params[":app"]
	version = this.GetString("version")
	date = this.GetString("date")
	channel = this.GetString("channel")
	tpage := this.GetString("page")
	page, err := strconv.Atoi(tpage)
	if err != nil {
		page = 1
	}
	dbcrash = http.GetFilteredCrashObj(app, version, channel, date, "")

	mapDate, mapType, mapVersion, mapChannel, mapDevice, mapOs = make(map[string]int), make(map[string]int),
		make(map[string]int), make(map[string]int), make(map[string]int), make(map[string]int)
	var mapCrash = make(map[string]CrashLog)
	for _, c := range dbcrash {
		m, ok := mapCrash[c.Md5]
		if ok {
			m.Count++
		} else {
			mapCrash[c.Md5] = CrashLog{c.Crashtype, "test", c.Md5, 1}
		}
		mapDate[c.Date]++
		mapType[c.Crashtype]++
		mapVersion[c.App]++
		mapChannel[c.Ch]++
		mapDevice[c.Dm]++
		mapOs[c.Os]++
	}
	for _, v := range mapCrash {
		crashLog = append(crashLog, v)
	}

	for k, _ := range mapDate {
		allDate = append(allDate, k)
	}
	for k, _ := range mapType {
		allType = append(allType, k)
	}
	for k, _ := range mapVersion {
		allVersion = append(allVersion, k)
	}
	for k, _ := range mapChannel {
		allChannel = append(allChannel, k)
	}
	for k, _ := range mapDevice {
		allDevice = append(allDevice, k)
	}
	for k, _ := range mapOs {
		allOs = append(allOs, k)
	}
	this.TplNames = "app.html"
	this.Data["App"] = app
	this.Data["Total"] = len(dbcrash)
	this.Data["AllCrashLog"] = crashLog
	this.Data["MapDate"] = mapDate
	this.Data["MapVersion"] = mapVersion
	this.Data["MapChannel"] = mapChannel
	this.Data["MapType"] = mapType
	this.Data["AllDate"] = allDate
	this.Data["AllVersion"] = allVersion
	this.Data["AllChannel"] = allChannel
	this.Data["CurPage"] = page
	this.Data["TotalPage"] = len(crashLog)/20 + 1

	//if app == tapp && version == tversion && date == tdate && channel == tchannel && crashType == tcrashType {
	//	allVersion, _ = this.GetSession("allVersion").(string)
	//	allDate, _ = this.GetSession("allDate").(string)
	//	allChannel, _ = this.GetSession("allChannel").(string)

	//	mapVersion, _ = this.GetSession("mapVersion").(map[string]int)
	//	mapDate, _ = this.GetSession("mapDate").(map[string]int)
	//	mapChannel, _ = this.GetSession("mapChannel").(map[string]int)
	//	mapType, _ = this.GetSession("mapType").(map[string]int)

	//	total, _ = this.GetSession("total").(int)
	//	crashLog, _ = this.GetSession("crashLog").([]CrashLog)
	//	crashCount = this.GetSession("crashCount").(map[string][]http.CrashObj)

	//} else {
	//	crashObj, tmapVersion, tmapChannel, tmapDate := http.GetFilteredCrashObj(tapp, tversion, tchannel, tdate, tcrashType)
	//	mapVersion = tmapVersion
	//	mapChannel = tmapChannel
	//	mapDate = tmapDate
	//	//解析所有可选条件
	//	if app != tapp {
	//		allVersion = ""
	//		allChannel = ""
	//		allDate = ""
	//		for k, _ := range mapVersion {
	//			allVersion = allVersion + "\n" + k
	//		}
	//		for k, _ := range mapChannel {
	//			allChannel = allChannel + "\n" + k
	//		}
	//		for k, _ := range mapDate {
	//			allDate = allDate + "\n" + k
	//		}
	//	} else {
	//		allVersion, _ = this.GetSession("allVersion").(string)
	//		allDate, _ = this.GetSession("allDate").(string)
	//		allChannel, _ = this.GetSession("allChannel").(string)
	//	}
	//	//记录下数据
	//	this.SetSession("app", tapp)
	//	this.SetSession("version", tversion)
	//	this.SetSession("date", tdate)
	//	this.SetSession("channel", tchannel)
	//	this.SetSession("crashType", tcrashType)

	//	this.SetSession("allVersion", allVersion)
	//	this.SetSession("allDate", allDate)
	//	this.SetSession("allChannel", allChannel)

	//	this.SetSession("mapVersion", mapVersion)
	//	this.SetSession("mapDate", mapDate)
	//	this.SetSession("mapChannel", mapChannel)

	//	total = len(crashObj)
	//	//crash log对应的crash obj
	//	crashCount = make(map[string][]http.CrashObj)
	//	this.SetSession("total", total)

	//	for _, v := range crashObj {
	//		var log string
	//		//取Log前300字符作为description
	//		if len(v.Log) >= 300 {
	//			log = v.Log[:300]
	//		} else {
	//			log = v.Log
	//		}
	//		crashCount[log] = append(crashCount[log], v)
	//	}
	//	this.SetSession("crashCount", crashCount)
	//	crashLog = make([]CrashLog, len(crashCount))
	//	index := 0
	//	mapType = make(map[string]int)
	//	re, _ := regexp.Compile("^\\w+\\.\\w+\\.\\w+")
	//	for log, count := range crashCount {
	//		name := re.FindString(log)
	//		if name == "" {
	//			name = log[:50]
	//		}
	//		crashLog[index] = CrashLog{name, log, len(count)}
	//		mapType[name] += len(count)
	//		index++
	//	}
	//	this.SetSession("mapType", mapType)
	//	crashLog = crashLog[:index]
	//	for i := 0; i < index; i++ {
	//		max := i
	//		for j := i + 1; j < index; j++ {
	//			if crashLog[j].Count > crashLog[max].Count {
	//				max = j
	//			}
	//		}
	//		if max != i {
	//			tmp := crashLog[i]
	//			crashLog[i] = crashLog[max]
	//			crashLog[max] = tmp
	//		}
	//	}
	//	this.SetSession("crashLog", crashLog)
	//}
	//this.SetSession("page", page)

	//this.TplNames = "app.html"
	//this.Data["App"] = tapp
	//this.Data["Total"] = total

	//tmp := strings.Split(allVersion, "\n")[1:]
	//sort.Strings(tmp)
	//this.Data["AllVersion"] = tmp
	//this.Data["MapVersion"] = mapVersion

	//tmp = strings.Split(allDate, "\n")[1:]
	//sort.Strings(tmp)
	//this.Data["AllDate"] = tmp
	//this.Data["MapDate"] = mapDate

	//tmp = strings.Split(allChannel, "\n")[1:]
	//sort.Strings(tmp)
	//this.Data["AllChannel"] = tmp
	//this.Data["MapChannel"] = mapChannel

	//this.Data["MapType"] = mapType

	//var totalPage int
	//totalPage = len(crashLog) / CRASH_PER_PAGE
	//if len(crashLog)%CRASH_PER_PAGE != 0 {
	//	totalPage++
	//}
	//this.Data["TotalPage"] = totalPage
	//this.Data["CurPage"] = page
	//if (page-1)*CRASH_PER_PAGE < len(crashLog) {
	//	if (page-1)*CRASH_PER_PAGE+CRASH_PER_PAGE < len(crashLog) {
	//		this.Data["AllCrashLog"] = crashLog[(page-1)*CRASH_PER_PAGE : (page-1)*CRASH_PER_PAGE+CRASH_PER_PAGE]
	//	} else {
	//		this.Data["AllCrashLog"] = crashLog[(page-1)*CRASH_PER_PAGE : len(crashLog)]
	//	}
	//}
}
