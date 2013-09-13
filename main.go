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
		for true {
			t := time.Now()
			n := time.Date(t.Year(), t.Month(), t.Day()+1, 6, 0, 0, 0, t.Location())
			duration := time.Duration(n.Unix()-t.Unix()) * time.Second
			fmt.Println("duration:", duration)
			<-time.After(duration)
			initialData()
		}
	}()

}

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crash/:crash", &controllers.CrashController{})
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

}
