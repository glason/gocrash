package main

import (
	"crash.android.meituan/controllers"
	"crash.android.meituan/models"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

func init() {
	go func() {
		initialData()
		t := time.NewTicker(time.Hour)
		for _ = range t.C {
			http.PeriodTask()
		}
	}()

}

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crash", &controllers.CrashController{})
	beego.Router("/:app", &controllers.AppController{})
	beego.Run()
}

func initialData() {
	fmt.Println("initial data time:", time.Now())
	if err := http.InitialCrashData(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("initial data successful!")
	}
	http.PeriodTask()
}
