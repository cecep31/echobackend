package main

import (
	"echobackend/config"
	"echobackend/internal/di"
	"runtime"
)

func main() {
	// Set GOMAXPROCS to match available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// load config
	conf, errconf := config.Load()
	if errconf != nil {
		panic(errconf)
	}

	// Initialize and run the application with wire
	app, err := di.BuildApplication(conf)
	if err != nil {
		panic(err)
	}

	// Run the application
	app.Run()
}
