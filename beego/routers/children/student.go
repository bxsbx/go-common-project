package children

import (
	"StandardProject/beego/controllers"
	"StandardProject/beego/routers/auth"
	"github.com/astaxie/beego"
)

func Student() {

	beego.InsertFilter("/v1/student/*", beego.BeforeRouter, auth.AuthStudentToken)

	beego.Router("/v1/student/web/test/test1", &controllers.TestController{}, "get:Test1")

	beego.Router("/v1/student/web/test/test2", &controllers.TestController{}, "get:Test2")

	// router general tag
}
