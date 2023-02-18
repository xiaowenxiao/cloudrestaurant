package main

import (
	"cloudrestaurant/controller"
	"cloudrestaurant/tool"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func main() {
	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}
	_, err = tool.OrmEngine(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	app := gin.Default()
	registerRouter(app)
	app.Run(cfg.AppHost + ":" + cfg.AppPort)
}

// 路由注册
func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
	new(controller.MemberController).Router(router)
}
