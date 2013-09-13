package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//崩溃日志jason对应的struct
type CrashObj struct {
	Uid   interface{}
	App   string //版本号
	Os    string //系统版本
	Appnm string //应用名称movie/hotel/xm/group
	Sc    string //屏幕大小
	Did   string
	Net   string //网络
	Ct    string //?
	City  string //城市
	Evs   []CrashEvs
	Dm    string //设备型号
	Uuid  string
	Ch    string //渠道
	Log   string //崩溃日志
	Date  string //日期
}

type CrashEvs struct {
	Nm  string
	Val CrashEvsVal
}

type CrashEvsVal struct {
	Log string
}

//crash服务器
const CRASH_URL = "http://10.64.11.188:8080/crash/"

//获取近7天的数据
const CRASH_DAYS = 7

var AllCrashObj []CrashObj
var crashDate string

func init() {
	AllCrashObj = make([]CrashObj, 0)
}

//在服务器http://10.64.11.188:8080/crash/上解析所有崩溃日志文件
func InitialCrashData() error {
	t := time.Now()
	for i := 1; i <= CRASH_DAYS; i++ {
		t = t.AddDate(0, 0, -1)
		y, m, d := t.Date()
		date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
		if !strings.Contains(crashDate, date) {
			fmt.Println("******getting data on ", date, " ********")
			crashDate = crashDate + date + "\n"
			if err := GetAllJsonObject(date); err != nil {
				fmt.Println(err)
			}
		}
	}
	//往前第八天的数据删除
	t = t.AddDate(0, 0, -1)
	y, m, d := t.Date()
	date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
	if strings.Contains(crashDate, date) {
		crashDate = strings.Replace(crashDate, date+"\n", "", 1)
		tmpCrashObj := make([]CrashObj, len(AllCrashObj))
		var index int
		for _, v := range AllCrashObj {
			if v.Date != date {
				//AllCrashObj = append(AllCrashObj[:k], AllCrashObj[k+1:]...)
				tmpCrashObj[index] = v
				index += 1
			}
		}
		AllCrashObj = tmpCrashObj[:index]
	}
	return nil
}

//解析崩溃object
func GetAllJsonObject(date string) error {
	url := CRASH_URL + "crash-android-" + date + ".json"
	//最多尝试3次
	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
	}
	if resp == nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	re, _ := regexp.Compile("\\{.*\\}")
	allJson := re.FindAll(body, -1)

	tmp := make([]CrashObj, len(allJson))

	re, _ = regexp.Compile("(\\n+|\\s{2,})")

	for i, v := range allJson {
		json.Unmarshal(v, &tmp[i])

		for _, v := range tmp[i].Evs {
			if v.Val.Log != "" {
				tmp[i].Log = re.ReplaceAllString(v.Val.Log, " ")
				break
			}
		}
		tmp[i].Date = date
	}
	AllCrashObj = append(AllCrashObj, tmp...)
	return nil
}

func GetFilteredCrashObj(app, version, channel, date, crashType string) ([]CrashObj, map[string]int, map[string]int, map[string]int) {
	result := make([]CrashObj, 0)
	tmp := make([]CrashObj, 100)
	index := 0
	//报表数据统计
	mapVersion := make(map[string]int)
	mapChannel := make(map[string]int)
	mapDate := make(map[string]int)
	for _, v := range AllCrashObj {
		if (app == "" || app == v.Appnm) && (version == "" || version == v.App) && (channel == "" || channel == v.Ch) && (date == "" || date == v.Date) && (crashType == "" || strings.Contains(v.Log, crashType)) {
			tmp[index] = v
			index++
			if index >= 100 {
				result = append(result, tmp...)
				index = 0
			}
			mapVersion[v.App] += 1
			mapChannel[v.Ch] += 1
			mapDate[v.Date] += 1
		}
	}
	if index < 100 {
		result = append(result, tmp[:index]...)
	}
	return result, mapVersion, mapChannel, mapDate
}
