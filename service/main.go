package main

import (
	"flag"
	"log"
	"sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/router"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("c", "conf.ini", "Path to configuration file")
	flag.Parse()

	// Initialize the application with the specified config path
	err := initialize.InitApp(*configPath)
	if err != nil {
		log.Panicln("初始化错误:", err)
	}

	httpPort := global.Config.GetValueString("base", "http_port")
	if err := router.InitRouters(":" + httpPort); err != nil {
		panic(err)
	}
}
