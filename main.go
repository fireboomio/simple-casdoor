package main

import (
	"casdoor/conf"
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
	beego.SetStaticPath("/swagger", "swagger")

	beego.InsertFilter("*", beego.BeforeRouter, routers.AutoSigninFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.CorsFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.AuthzFilter)
	beego.InsertFilter("*", beego.BeforeRouter, routers.RecordMessage)

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "casdoor_session_id"
	if conf.GetConfigString("redisEndpoint") == "" {
		beego.BConfig.WebConfig.Session.SessionProvider = "file"
		beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	} else {
		beego.BConfig.WebConfig.Session.SessionProvider = "redis"
		beego.BConfig.WebConfig.Session.SessionProviderConfig = conf.GetConfigString("redisEndpoint")
	}
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24 * 30

	err := logs.SetLogger(logs.AdapterFile, conf.GetConfigString("logConfig"))
	if err != nil {
		panic(err)
	}
	port := beego.AppConfig.DefaultInt("httpport", 10021)
	// logs.SetLevel(logs.LevelInformational)
	logs.SetLogFuncCall(false)

	beego.Run(fmt.Sprintf(":%v", port))
}
