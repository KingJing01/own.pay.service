package main

import (
	_ "paymentservice/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	var FilterUser = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json

		// token := ctx.Request.Header.Get("Authorization")
		// if token == "" && ctx.Request.RequestURI != "/v1/authoritymanage/AuthorityError" && ctx.Request.RequestURI != "/v1/authoritymanage/Login" {
		// 	ctx.Redirect(302, "/v1/authoritymanage/AuthorityError")

		// } else if token != "" && ctx.Request.RequestURI != "/v1/authoritymanage/AuthorityError" && ctx.Request.RequestURI != "/v1/authoritymanage/Login" {
		// 	result, _, _ := tools.CheckLogin(token)
		// 	if result == false {
		// 		ctx.Redirect(302, "/v1/authoritymanage/AuthorityError")

		// 	}
		// }

		//fmt.Println("Number of records deleted in database:", username)
	}

	beego.InsertFilter("*", beego.BeforeRouter, FilterUser)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
