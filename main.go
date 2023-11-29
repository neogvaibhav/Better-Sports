package main

import (
	"github.com/Basu008/Better-ESPN/app"
	"github.com/Basu008/Better-ESPN/config"
)

func main() {
	config := config.NewConfig()
	app.SetUpApp(config)
}
