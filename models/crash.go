package crash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//崩溃日志jason对应的struct
type CrashObj struct {
	Uid   interface{}
	App   string
	Os    string
	Appnm string
	Sc    string
	Did   string
	Net   string
	Ct    string
	City  string
	Evs   []CrashEvs
	Dm    string
	Uuid  string
	Ch    string
	Log   string
	Date  string
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

func init() {
	AllCrashObj = make([]CrashObj, 0)
}

//在服务器http://10.64.11.188:8080/crash/上解析所有崩溃日志文件
func InitialCrashData() error {
	fmt.Println("get all crash data url...")
	resp, err := http.Get(CRASH_URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	src := string(body)

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")

	re, _ = regexp.Compile("crash(\\d||-)*\\.json")
	allJson := re.FindAllString(src, -1)

	for _, v := range allJson {
		fmt.Println("get all crash obj in " + v)
	}
	if err := GetAllJsonObject("2013-08-29"); err != nil {
		return err
	}
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

func GetFilteredCrashObj(app, version, channel, date string) []CrashObj {
	result := make([]CrashObj, 0)
	tmp := make([]CrashObj, 100)
	index := 0
	for _, v := range AllCrashObj {
		if (app == "" || app == v.Appnm) && (version == "" || version == v.App) && (channel == "" || channel == v.Ch) && (date == "" || date == v.Date) {
			tmp[index] = v
			index++
			if index >= 100 {
				result = append(result, tmp...)
				index = 0
			}
		}
	}
	if index < 100 {
		fmt.Println("append")
		result = append(result, tmp[:index]...)
	}
	return result
}
