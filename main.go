package main

import (
	"ledger-app/config"
	"ledger-app/internal/connections/echoserver"
	"ledger-app/internal/providers"
)

func main() {
	e := echoserver.GetInstance()
	cfg := config.LoadEnvironment()

	providers.InitLogger()
	providers.InitDatabase()
	providers.RegisterMiddlewares(e)
	providers.InitDefaultAdmin()
	providers.StartServer(e, cfg)
}
