package main

import (
	"os"
	"os/signal"
	"syscall"
	"vega-mm/api"
	"vega-mm/bot"
	"vega-mm/logging"
	"vega-mm/store"
	"vega-mm/vega"
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
