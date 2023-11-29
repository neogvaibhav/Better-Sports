package main

import (
	"github.com/neogvaibhav/Better-Sports/app"
	"github.com/neogvaibhav/Better-Sports/config"
)

func main() {
	config := config.NewConfig()
	app.SetUpApp(config)
}
