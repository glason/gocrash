package main

import (
	"crash.android.meituan/controllers"
	"crash.android.meituan/models"
	"fmt"
	"github.com/astaxie/beego"
)

func init() {
	if err := http.InitialCrashData(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("initial data successful!")
	}

}

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crash/:crash", &controllers.CrashController{})
	beego.Router("/:app", &controllers.AppController{})
	beego.Run()
}
