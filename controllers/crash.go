package controllers

import (
	"crash.android.meituan/models"
	"github.com/astaxie/beego"
	"regexp"

	////"net/url"
	"strconv"

//"strings"
)

type CrashController struct {
	beego.Controller
}

func (this *CrashController) Get() {
	var app, date, version, channel, md5 string
	var count int
	var dbCrash []http.Dbcrash
	var dateStats, versionStats, deviceStats, osStats []CrashStats

	app = this.GetString("app")
	date = this.GetString("date")
	version = this.GetString("version")
	channel = this.GetString("channel")
	md5 = this.GetString("md5")

	var tmpmap [][]map[string][]byte

	count, dbCrash, tmpmap = http.GetDataForCrashPage(app, date, version, channel, md5)

	for _, m := range tmpmap[0] {
		s := string(m["date"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		dateStats = append(dateStats, t)
	}
	for _, m := range tmpmap[1] {
		s := string(m["app"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		versionStats = append(versionStats, t)
	}
	for _, m := range tmpmap[2] {
		s := string(m["dm"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		deviceStats = append(deviceStats, t)
	}
	for _, m := range tmpmap[3] {
		s := string(m["os"])
		c, _ := strconv.Atoi(string(m["count(*)"]))
		t := CrashStats{s, c}
		osStats = append(osStats, t)
	}

	detail := dbCrash[0].Log
	re, _ := regexp.Compile("java:\\d+")
	detail = re.ReplaceAllString(detail, "java:***")

	this.TplNames = "crash.html"
	this.Data["Appnm"] = app
	this.Data["Total"] = count
	this.Data["CrashType"] = dbCrash[0].Crashtype
	this.Data["CrashDetail"] = detail
	this.Data["Dbcrash"] = dbCrash
	this.Data["DateStats"] = dateStats
	this.Data["VersionStats"] = versionStats
	this.Data["DeviceStats"] = deviceStats
	this.Data["OsStats"] = osStats

	//this.Data["Appnm"], _ = this.GetSession("app").(string)

	//crashID := this.Ctx.Params[":crash"]
	//index, err := strconv.Atoi(crashID)
	//if err != nil {
	//	return
	//}
	//var crashObj []http.CrashObj
	//crashLog, _ := this.GetSession("crashLog").([]CrashLog)
	//crashCount, _ := this.GetSession("crashCount").(map[string][]http.CrashObj)
	//page, _ := this.GetSession("page").(int)

	//i := index + (page-1)*CRASH_PER_PAGE
	//if i < len(crashLog) {
	//	log := crashLog[i].Description
	//	crashObj = crashCount[log]
	//}
	//if crashObj == nil {
	//	return
	//}
	//this.Data["Total"] = len(crashObj)
	//this.Data["CrashObj"] = crashObj
	//if len(crashObj) >= 1 {
	//	log := crashObj[0].Log
	//	if strings.Index(log, ":") > 0 {
	//		this.Data["CrashType"] = log[:strings.Index(log, ":")]
	//	}
	//	this.Data["CrashDetail"] = log
	//}
	////报表
	//mapVersion := make(map[string]int)
	//mapDate := make(map[string]int)
	//mapChannel := make(map[string]int)
	//mapCity := make(map[string]int)
	//mapDevice := make(map[string]int)
	//mapScreen := make(map[string]int)
	//mapOS := make(map[string]int)
	//mapNet := make(map[string]int)
	//for _, v := range crashObj {
	//	mapVersion[v.App] += 1
	//	mapChannel[v.Ch] += 1
	//	mapDate[v.Date] += 1
	//	mapCity[v.City] += 1
	//	mapDevice[v.Dm] += 1
	//	mapScreen[v.Sc] += 1
	//	mapOS[v.Os] += 1
	//	mapNet[v.Net] += 1
	//}
	//this.Data["MapVersion"] = mapVersion
	//this.Data["MapDate"] = mapDate
	//this.Data["MapChannel"] = mapChannel
	//this.Data["MapCity"] = mapCity
	//this.Data["MapDevice"] = mapDevice
	//this.Data["MapScreen"] = mapScreen
	//this.Data["MapOS"] = mapOS
	//this.Data["MapNet"] = mapNet

}
