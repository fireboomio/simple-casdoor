package routers

import (
	"github.com/beego/beego"
	"github.com/beego/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["casdoor/controllers:ApiController"] = append(beego.GlobalControllerRouter["casdoor/controllers:ApiController"],
        beego.ControllerComments{
            Method: "AddUser",
            Router: `/add-user`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["casdoor/controllers:ApiController"] = append(beego.GlobalControllerRouter["casdoor/controllers:ApiController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["casdoor/controllers:ApiController"] = append(beego.GlobalControllerRouter["casdoor/controllers:ApiController"],
        beego.ControllerComments{
            Method: "SendVerificationCode",
            Router: `/send-verification-code`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["casdoor/controllers:ApiController"] = append(beego.GlobalControllerRouter["casdoor/controllers:ApiController"],
        beego.ControllerComments{
            Method: "UpdateProvider",
            Router: `/update-provider`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
