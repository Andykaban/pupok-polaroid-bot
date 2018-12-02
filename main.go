package main

import (
	"github.com/Andykaban/pupok-polaroid-bot/transform"
	"log"
)

func main() {
	log.Println("Start bot..")
	transformer, err := transform.New("./static/images/background.png",
		"./static/fonts/wqy-zenhei.ttf")
	if err != nil {
		log.Fatal(err)
	}
	err = transformer.CreatePolaroidImage("/Users/andy/Downloads/cat.jpg",
		"/Users/andy/Downloads/cat_1.jpg", "!!! KURLIKURLI !!!")
	if err != nil {
		log.Fatal(err)
	}
}
