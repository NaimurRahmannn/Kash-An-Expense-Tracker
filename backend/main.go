package main

import (
	"os"

	_ "backend/routers"
	"backend/storage"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := storage.Init(); err != nil {
		logs.Error("failed to initialize storage: %v", err)
		os.Exit(1)
	}
	defer storage.Close()

	beego.Run()
}
