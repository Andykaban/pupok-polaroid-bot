package main

import (
	"flag"
	"github.com/Andykaban/pupok-polaroid-bot/bot"
	"github.com/Andykaban/pupok-polaroid-bot/config"
	"log"
)

func Usage() {
	log.Println("!!! Pupok polaroid bot !!!")
	log.Println("Usage: pupok-polaroid-bot -config 'path to config bot file'")
	flag.PrintDefaults()
}

func main() {
	log.Println("Start bot..")
	configPath := flag.String("config", "config.json", "path to bot config file")
	botConfig := config.ParseBotConfig(*configPath)
	tgBot := bot.NewBot(botConfig)
	tgBot.WatchDog()
	tgBot.BotMainHandler()
}
