package main

import (
	"os"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	_ "backend/routers"
	"backend/storage"
)

func main() {
	if port := os.Getenv("PORT"); port != "" {
		if parsedPort, err := strconv.Atoi(port); err == nil {
			beego.BConfig.Listen.HTTPPort = parsedPort
		}
	}

	if err := storage.Init(); err != nil {
		logs.Error("failed to initialize storage: %v", err)
		os.Exit(1)
	}
	defer storage.Close()

	beego.Run()
}
