package main

import (
	"github.com/HardDie/blog_engine/internal/application"
	"github.com/HardDie/blog_engine/internal/logger"
)

func main() {
	app, err := application.Get()
	if err != nil {
		logger.Error.Fatal(err)
	}
	logger.Info.Println("Server listen on", app.Cfg.Port)
	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err)
	}
}
