package main

import (
	"embed"
	"log"
	"sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/lib/embedfs"
	"sun-panel/router"
)

//go:embed assets/*
var embeddedAssets embed.FS

func main() {
	embedfs.Init(embeddedAssets)
	err := initialize.InitApp()
	if err != nil {
		log.Println("初始化错误:", err.Error())
		panic(err)
	}

	httpPort := global.Config.GetValueString("base", "http_port")
	if err := router.InitRouters(":" + httpPort); err != nil {
		panic(err)
	}
}
