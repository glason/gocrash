package main

import (
	"crash.android.meituan/controllers"
	"crash.android.meituan/models"
	"fmt"
	"github.com/astaxie/beego"
)

func init() {
	if err := crash.InitialCrashData(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("initial data successful!")
	}

}

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/:app", &controllers.CrashController{})
	beego.Run()
}
