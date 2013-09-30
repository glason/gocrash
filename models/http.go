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
CREATE TABLE `test`.`dbcrash` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `uid` VARCHAR(100) NULL,
  `app` VARCHAR(100) NULL,
  `os` VARCHAR(100) NULL,
  `appnm` VARCHAR(100) NULL,
  `sc` VARCHAR(100) NULL,
  `did` VARCHAR(100) NULL,
  `net` VARCHAR(100) NULL,
  `ct` VARCHAR(100) NULL,
  `city` VARCHAR(100) NULL,
  `dm` VARCHAR(100) NULL,
  `uuid` VARCHAR(100) NULL,
  `ch` VARCHAR(100) NULL,
  `log` TEXT NULL,
  `md5` VARCHAR(100) NULL,
  `date` VARCHAR(100) NULL,
  PRIMARY KEY (`id`));
*/

//数据库存储crash结构
type Dbcrash struct {
	Id        int `beedb:"PK"`
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

//获取近7天的数据
const CRASH_DAYS = 2

var crashDate string
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
	return nil
	t := time.Now()
	for i := 1; i <= CRASH_DAYS; i++ {
		t = t.AddDate(0, 0, -1)
		y, m, d := t.Date()
		date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
		if !strings.Contains(crashDate, date) {
			crashDate = crashDate + date + "\n"
			if err := GetAllJsonObject(date); err != nil {
				fmt.Println("******getting data on ", date, " failed********")
				fmt.Println(err)
			} else {
				fmt.Println("******getting data on ", date, " successful********")
			}
		}
	}
	//往前第八天的数据删除
	t = t.AddDate(0, 0, -1)
	y, m, d := t.Date()
	date := fmt.Sprintf("%.4d-%.2d-%.2d", y, m, d)
	if strings.Contains(crashDate, date) {
		crashDate = strings.Replace(crashDate, date+"\n", "", 1)
		orm.SetTable("dbcrash").Where("date=", date).DeleteRow()
	}
	return nil
}

//解析崩溃object
func GetAllJsonObject(date string) error {
	url := CRASH_URL + "crash-android-" + date + ".json"
	fmt.Println("start getting data on ", url)
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

	re, _ = regexp.Compile("(\\n+|\\s{2,}|java:\\d+)")
	var tmp CrashObj
	var log, _uid, crashType string
	for _, v := range allJson {
		json.Unmarshal(v, &tmp)
		for _, v := range tmp.Evs {
			if v.Val.Log != "" {
				log = v.Val.Log
				break
			}
		}
		h := md5.New()
		io.WriteString(h, re.ReplaceAllString(log, " "))
		_md5 := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println(_md5)
		if value, ok := tmp.Uid.(string); ok {
			_uid = value
		} else if value, ok := tmp.Uid.(int); ok {
			_uid = strconv.Itoa(value)
		} else {
			_uid = ""
		}
		crashType = getCrashType(log)

		dbCrash := Dbcrash{0, _uid, tmp.App, tmp.Os, tmp.Appnm, tmp.Sc, tmp.Did, tmp.Net, tmp.Ct, tmp.City, tmp.Dm, tmp.Uuid, tmp.Ch, log, _md5, date, crashType}

		fmt.Println(dbCrash)
		if err := orm.Save(&dbCrash); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func getCrashType(crash string) string {
	index := strings.Index(crash, "Caused by:")
	var subCrash string
	if index == -1 {
		subCrash = crash
	} else {
		subCrash = crash[index+len("Caused by:"):]
	}
	index = strings.Index(subCrash, ":")
	var crashType string
	if index == -1 {
		crashType = "TypeNotFound"
	} else {
		crashType = subCrash[:index]
	}
	return crashType
}

func GetFilteredCrashObj(app, version, channel, date, _md5 string) []Dbcrash {
	var filter string
	var result []Dbcrash
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
	if _md5 != "" {
		filter = filter + "and md5='" + _md5 + "' "
	}
	//if m, err := orm.SetTable("dbcrash").Where(filter).Select("id").FindMap(); err == nil {
	//	count = len(m)
	//}
	fmt.Println("filter:", filter)
	orm.Where(filter).Limit(0, 20).FindAll(&result)
	fmt.Println("*****GetFilteredCrashObj*****")
	fmt.Println(result)
	return result
}

//func GetCrashGroupByMd5(app, version, channel, date string) []map[string][]byte {
//	var filter string
//	var result []map[string][]byte
//	if app != "" {
//		filter = filter + "appnm='" + app + "' "
//	}
//	if version != "" {
//		filter = filter + "and app='" + version + "' "
//	}
//	if channel != "" {
//		filter = filter + "and ch='" + channel + "' "
//	}
//	if date != "" {mysql.gomysql.gomysql.go
//		filter = filter + "and date='" + date + "' "
//	}
//	result, _ = orm.Select("log,count(*)").Where(filter).GroupBy("md5").OrderBy("count(*) desc").FindMap()
//	fmt.Println("*****GetCrashGroupByMd5*****")
//	fmt.Println(result)
//	return result
//}

//func GetCrashMap(filter, group string) []map[string][]byte {
//	var result map[string]int
//	result = make(map[string]int)
//	dbMap := orm.SetTable("dbcrash").Select(group + ",count(*)").Where(filter).GroupBy(group).OrderBy("count(*) desc").FindMap()
//	for _, m := range dbMap {
//		key := string(m[group])
//		value := int()
//	}
//}
