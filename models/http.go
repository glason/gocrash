package http

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/ziutek/mymysql/godrv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
}

type CrashEvs struct {
	Nm  string
	Val CrashEvsVal
}

type CrashEvsVal struct {
	Log string
}

/*
CREATE TABLE `dbcrash` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` varchar(100) DEFAULT NULL,
  `app` varchar(100) DEFAULT NULL,
  `os` varchar(100) DEFAULT NULL,
  `appnm` varchar(100) DEFAULT NULL,
  `sc` varchar(100) DEFAULT NULL,
  `did` varchar(100) DEFAULT NULL,
  `net` varchar(100) DEFAULT NULL,
  `ct` varchar(100) DEFAULT NULL,
  `city` varchar(100) CHARACTER SET utf8 DEFAULT NULL,
  `dm` varchar(100) DEFAULT NULL,
  `uuid` varchar(100) DEFAULT NULL,
  `ch` varchar(100) DEFAULT NULL,
  `log` text,
  `md5` varchar(100) DEFAULT NULL,
  `date` varchar(100) DEFAULT NULL,
  `crashtype` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ;

*/

//数据库存储crash结构
type Dbcrash struct {
	Id        int `beedb:"PK"` //it is important to set PK
	Uid       string
	App       string
	Os        string
	Appnm     string
	Sc        string
	Did       string
	Net       string
	Ct        string
	City      string
	Dm        string
	Uuid      string
	Ch        string
	Log       string
	Md5       string //crash log md5，用于判断相同crash
	Date      string
	Crashtype string //crash类型
}

//crash服务器
const CRASH_URL = "http://10.64.11.188:8080/crash/"

var CRASH_URL_TODAY = []string{"http://10.64.12.213:8080/logs/clientCrash.json", "http://10.64.13.226:8080/logs/clientCrash.json",
	"http://10.64.11.188:8080/logs/clientCrash.json", "http://10.64.11.187:8080/logs/clientCrash.json"}

//const CRASH_URL = "http://127.0.0.1:8000/"

//获取近6天的数据
const CRASH_DAYS = 6

var UpdateTime time.Time
var orm beedb.Model

func init() {
	db, err := sql.Open("mymysql", "test/wangjiasheng/wangjiasheng")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
}

//在服务器http://10.64.11.188:8080/crash/上解析所有崩溃日志文件
func InitialCrashData() error {
	orm.SetTable("dbcrash").DeleteRow()
	//return nil
	t := time.Now()
	for i := 1; i <= CRASH_DAYS; i++ {
		t = t.AddDate(0, 0, -1)
		date := getDateString(t)

		if err := GetAllJsonObject(CRASH_URL+"crash-android-"+date+".json", date); err != nil {
			fmt.Println("******getting data on ", date, " failed********")
			fmt.Println(err)
		} else {
			fmt.Println("******getting data on ", date, " successful********")
		}
	}
	return nil
}

//每小时执行一次
func PeriodTask() {
	UpdateTime = time.Now()
	cmd := exec.Command("sh", "crash.sh")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	//删除当天数据以及第7天前数据
	t := []time.Time{time.Now(), time.Now().AddDate(0, 0, -7)}
	for _, v := range t {
		orm.SetTable("dbcrash").Where("date=?", getDateString(v)).DeleteRow()
	}
	GetAllJsonObject("", getDateString(time.Now()))

}

func getDateString(t time.Time) string {
	y, m, d := t.Date()
	date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
	return date
}

//解析崩溃object
func GetAllJsonObject(url, date string) error {
	fmt.Println("start getting data on ", url)
	//最多尝试3次
	var buf []byte
	var err error
	if url == "" {
		cmd := exec.Command("cat", "crash.json")
		buf, err = cmd.Output()
		if err != nil {
			return err
		}
	} else {
		var resp *http.Response
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

		buf, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	re, _ := regexp.Compile("\\{.*\\}")
	allJson := re.FindAll(buf, -1)

	//re, _ = regexp.Compile("(\\n+|\\s{2,}|java:\\d+)")
	re, _ = regexp.Compile("(\\s+|\\d+)")

	var log, _uid, crashType string
	var dbcount int
	dbcount = 0
	for _, v := range allJson {
		var tmp CrashObj
		json.Unmarshal(v, &tmp)
		if tmp.Ct != "android" {
			continue
		}
		for _, v := range tmp.Evs {
			if v.Val.Log != "" {
				log = v.Val.Log
				break
			}
		}
		h := md5.New()
		io.WriteString(h, re.ReplaceAllString(log, ""))
		_md5 := fmt.Sprintf("%x", h.Sum(nil))
		if value, ok := tmp.Uid.(string); ok {
			_uid = value
		} else if value, ok := tmp.Uid.(int); ok {
			_uid = strconv.Itoa(value)
		} else {
			_uid = ""
		}
		crashType = getCrashType(log)

		dbCrash := Dbcrash{0, _uid, tmp.App, tmp.Os, tmp.Appnm, tmp.Sc, tmp.Did, tmp.Net, tmp.Ct, tmp.City, tmp.Dm, tmp.Uuid, tmp.Ch, log, _md5, date, crashType}

		//fmt.Println(dbCrash)
		if err := orm.Save(&dbCrash); err != nil {
			fmt.Println(err)
		} else {
			dbcount += 1
		}
	}
	fmt.Println("getting crash count:", dbcount)
	return nil
}

func getCrashType(crash string) string {
	var crashType string
	re, _ := regexp.Compile("Caused\\s*by:\\s*[^\\s:]+")
	subCrash := re.FindString(crash)
	if subCrash == "" {
		re, _ = regexp.Compile("^[^\\s:]+")
		crashType = re.FindString(crash)
	} else {
		index := strings.Index(subCrash, ":")
		crashType = subCrash[index+1:]
	}
	return strings.TrimSpace(crashType)
}

func GetDataForAppPage(app, date, version, channel string, start, limit int) (int, int, [][]map[string][]byte) {
	var filter string
	var total, logcount int
	var result [][]map[string][]byte
	if app != "" {
		filter = filter + "appnm='" + app + "' "
	}
	if version != "" {
		filter = filter + "and app='" + version + "' "
	}
	if channel != "" {
		filter = filter + "and ch='" + channel + "' "
	}
	if date != "" {
		filter = filter + "and date='" + date + "' "
	}
	m, _ := orm.SetTable("dbcrash").Select("count(*) as count").Where(filter).FindMap()
	total, _ = strconv.Atoi(string(m[0]["count"]))

	m, _ = orm.SetTable("dbcrash").Select("count(distinct md5) as count").Where(filter).FindMap()
	logcount, _ = strconv.Atoi(string(m[0]["count"]))

	m, _ = orm.SetTable("dbcrash").Select("crashtype,log,md5,count(*)").Where(filter).GroupBy("md5").OrderBy("count(*) desc").Limit(limit, start).FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("date,count(*)").Where(filter).GroupBy("date").OrderBy("date").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("crashtype,count(*)").Where(filter).GroupBy("crashtype").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("app,count(*)").Where(filter).GroupBy("app").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("dm,count(*)").Where(filter).GroupBy("dm").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("ch,count(*)").Where(filter).GroupBy("ch").OrderBy("count(*) desc").FindMap()
	result = append(result, m)
	return total, logcount, result
}

func GetDataForCrashPage(app, date, version, channel, md5 string) (int, []Dbcrash, [][]map[string][]byte) {
	var filter string
	var count int
	var crashObj []Dbcrash
	var result [][]map[string][]byte
	if app != "" {
		filter = filter + "appnm='" + app + "' "
	}
	if date != "" {
		filter = filter + "and date='" + date + "' "
	}
	if version != "" {
		filter = filter + "and app='" + version + "' "
	}
	if channel != "" {
		filter = filter + "and ch='" + channel + "' "
	}
	if md5 != "" {
		filter = filter + "and md5='" + md5 + "' "
	}
	m, _ := orm.SetTable("dbcrash").Select("count(*) as count").Where(filter).FindMap()
	count, _ = strconv.Atoi(string(m[0]["count"]))

	orm.Where(filter).Limit(20, 0).FindAll(&crashObj)

	m, _ = orm.SetTable("dbcrash").Select("date,count(*)").Where(filter).GroupBy("date").OrderBy("date").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("app,count(*)").Where(filter).GroupBy("app").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("dm,count(*)").Where(filter).GroupBy("dm").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	m, _ = orm.SetTable("dbcrash").Select("os,count(*)").Where(filter).GroupBy("os").OrderBy("count(*) desc").FindMap()
	result = append(result, m)

	return count, crashObj, result
}
