package main

import (
	"casdoor/object"
	"casdoor/routers"
	"flag"
	"fmt"
	"github.com/beego/beego"
	"github.com/beego/beego/logs"
)

// 命令行参数 -createDatabase
func getCreateDatabaseFlag() bool {
	res := flag.Bool("createDatabase", false, "true if you need Casdoor to create database")
	flag.Parse()
	return *res
}

func main() {
	createDatabase := getCreateDatabaseFlag()

	object.InitAdapter()
	object.CreateTables(createDatabase)
	object.InitDb()

	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.CopyRequestBody = true
	beego.SetStaticPath("/swagger", "swagger")

	beego.InsertFilter("*", beego.BeforeRouter, routers.AutoSigninFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.CorsFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.AuthzFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.RecordMessage)

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "casdoor_session_id"
	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 * 30

	err := logs.SetLogger(logs.AdapterFile, "{\"filename\": \"logs/casdoor.log\", \"maxdays\":99999, \"perm\":\"0770\"}")
	if err != nil {
		panic(err)
	}
	port := beego.AppConfig.DefaultInt("httpport", 10021)
	logs.SetLogFuncCall(false)

	beego.Run(fmt.Sprintf(":%v", port))
}
