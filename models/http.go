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

const CRASH_URL = "http://10.64.11.188:8080/crash/"

var AllCrashObj []CrashObj
var crashDate string

func init() {
	AllCrashObj = make([]CrashObj, 0)
}

//在服务器http://10.64.11.188:8080/crash/上解析所有崩溃日志文件
func InitialCrashData() error {
	fmt.Println("get all crash data url...")
	//resp, err := http.Get(CRASH_URL)
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//src := string(body)

	////去除所有尖括号内的HTML代码，并换成换行符
	//re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	//src = re.ReplaceAllString(src, "")

	//re, _ = regexp.Compile("crash(\\d||-)*\\.json")
	//allJson := re.FindAllString(src, -1)

	//for _, v := range allJson {
	//	fmt.Println("get all crash obj in " + v)
	//}

	t := time.Now()
	for i := 1; i <= 7; i++ {
		t = t.AddDate(0, 0, -1)
		y, m, d := t.Date()
		date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
		fmt.Println("date:", date)
		if !strings.Contains(crashDate, date) {
			crashDate = crashDate + date + "\n"
			if err := GetAllJsonObject(date); err != nil {
				fmt.Println(err)
			}
		}
	}
	t = t.AddDate(0, 0, -1)
	y, m, d := t.Date()
	date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
	if strings.Contains(crashDate, date) {
		crashDate = strings.Replace(crashDate, date+"\n", "", 1)
		for k, v := range AllCrashObj {
			if v.Date == date {
				AllCrashObj = append(AllCrashObj[:k], AllCrashObj[k+1:]...)
			}
		}
	}

	//if err := GetAllJsonObject("2013-08-29"); err != nil {
	//	fmt.Println(err)
	//}
	return nil
}

//解析崩溃object
func GetAllJsonObject(date string) error {
	fmt.Println("getting all json in " + date + "...")
	url := CRASH_URL + "crash-android-" + date + ".json"
	resp, err := http.Get(url)
	if err != nil {
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

func GetFilteredCrashObj(app, version, channel, date string) ([]CrashObj, map[string]int, map[string]int, map[string]int) {
	result := make([]CrashObj, 0)
	tmp := make([]CrashObj, 100)
	index := 0
	//报表数据统计
	mapVersion := make(map[string]int)
	mapChannel := make(map[string]int)
	mapDate := make(map[string]int)
	for _, v := range AllCrashObj {
		if (app == "" || app == v.Appnm) && (version == "" || version == v.App) && (channel == "" || channel == v.Ch) && (date == "" || date == v.Date) {
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
		fmt.Println("append")
		result = append(result, tmp[:index]...)
	}
	return result, mapVersion, mapChannel, mapDate
}
