package bot

import (
	"github.com/Andykaban/pupok-polaroid-bot/config"
	"github.com/Andykaban/pupok-polaroid-bot/transform"
	"github.com/Andykaban/pupok-polaroid-bot/utils"
	"github.com/disintegration/imaging"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"image/jpeg"
	"log"
	"net/http"
)

const TELEGRAM_ROOT = "https://api.telegram.org/file/bot"

type Bot struct {
	TelegramBot *tgbotapi.BotAPI
	Transformer *transform.PolaroidTransform
	TempDir string
	LabelsMap map[int]string
}

func NewBot(botConfig *config.BotConfig) *Bot {
	var bot *tgbotapi.BotAPI
	var err error
	if botConfig.BotProxyUrl == "" {
		bot, err = tgbotapi.NewBotAPI(botConfig.BotToken)
		if err != nil {
			panic(err)
		}
	} else {
		var auth proxy.Auth
		if botConfig.BotProxyLogin != "" || botConfig.BotProxyPassword != "" {
			auth = proxy.Auth{User:botConfig.BotProxyLogin, Password:botConfig.BotProxyPassword}
		}
		dialer, err := proxy.SOCKS5("tcp", botConfig.BotProxyUrl, &auth, proxy.Direct)
		if err != nil {
			panic(err)
		}
		httpTransport := &http.Transport{}
		httpTransport.Dial = dialer.Dial
		httpClient := http.Client{Transport:httpTransport}
		bot, err = tgbotapi.NewBotAPIWithClient(botConfig.BotToken, &httpClient)
		if err != nil {
			panic(err)
		}
		http.DefaultTransport = httpTransport
	}
	transformer, err := transform.New(botConfig.BackgroundPath, botConfig.FontPath)
	if err != nil {
		panic(err)
	}
	return &Bot{TelegramBot:bot,
		Transformer:transformer,
		TempDir:botConfig.BotTempDir,
		LabelsMap:make(map[int]string)}
}

func (b *Bot) sendTextMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	b.TelegramBot.Send(msg)
}

func (b *Bot) downloadPhoto(chatId int64, photoId, downloadPath string) {
	log.Printf("Try to download file with %s ID", photoId)
	resp, err := b.TelegramBot.GetFile(tgbotapi.FileConfig{FileID:photoId})
	if err != nil {
		log.Println(err)
		b.sendTextMessage(chatId, err.Error())
		return
	}
	downloadUrl := TELEGRAM_ROOT + b.TelegramBot.Token + "/" +resp.FilePath
	raw, err := http.Get(downloadUrl)
	defer raw.Body.Close()
	if err != nil {
		log.Println(err)
		b.sendTextMessage(chatId, err.Error())
		return
	}
	img, err := jpeg.Decode(raw.Body)
	if err != nil {
		log.Println(err)
		b.sendTextMessage(chatId, err.Error())
	}
	imaging.Save(img, downloadPath)
}

func (b *Bot) BotMainHandler() {
	log.Printf("Authorized on account %s", b.TelegramBot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := b.TelegramBot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatId := update.Message.Chat.ID
		if update.Message.Photo != nil {
			downloadFileName := utils.GetRandomFileName(b.TempDir)
			newFileName := utils.GetRandomFileName(b.TempDir)
			photos := *update.Message.Photo
			photoId := photos[1].FileID
			b.downloadPhoto(chatId, photoId, downloadFileName)
			b.Transformer.CreatePolaroidImage(downloadFileName, newFileName, "PAHOM")
		}
	}
}
