package main

import (
	"os"
	"os/signal"
	"syscall"
	"vega-cli-mm/api"
	"vega-cli-mm/bot"
	"vega-cli-mm/logging"
	"vega-cli-mm/store"
	"vega-cli-mm/vega"
)

const CoreNode = "darling.network:3007"

func keepAlive() {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	<-gracefulStop
	logging.GetLogger().Info("shutting down on user request")
}

func main() {
	appStore := store.NewStore()
	vegaClient := vega.NewVega(appStore, CoreNode)
	bot.NewBot(appStore, vegaClient).Start()
	api.NewApi(appStore).Start()
	keepAlive()
}
