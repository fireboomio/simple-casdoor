package routers

import (
	"casdoor/controllers"
	"github.com/beego/beego"
)

func init() {
	initAPI()
}

func initAPI() {

	beego.Router("/api/login", &controllers.ApiController{}, "POST:Login")

	beego.Router("/api/add-user", &controllers.ApiController{}, "POST:AddUser")

	beego.Router("/api/get-user", &controllers.ApiController{}, "GET:GetUserByToken")

	beego.Router("/api/send-verification-code", &controllers.ApiController{}, "POST:SendVerificationCode")

	beego.Router("/api/update-provider", &controllers.ApiController{}, "POST:UpdateProvider")
}
