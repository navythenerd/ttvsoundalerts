package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ttvsoundalerts/twitch"
)

func main() {
	log.Println("TTV Sound Alerts Bot")
	log.Println("by NavyTheNerd")

	var cfg twitch.Config
	twitch.ReadConfig(&cfg, "config.json")

	bot := twitch.New(&cfg)
	bot.Connect()

	defer bot.Shutdown()

	http.ListenAndServe(":8000", bot)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}
