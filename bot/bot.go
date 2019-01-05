package bot

import (
	"fmt"
	"github.com/Andykaban/pupok-polaroid-bot/config"
	"github.com/Andykaban/pupok-polaroid-bot/transform"
	"github.com/Andykaban/pupok-polaroid-bot/utils"
	"github.com/disintegration/imaging"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

const telegramRoot = "https://api.telegram.org/file/bot"
const botTaskCap = 30

type Bot struct {
	TelegramBot *tgbotapi.BotAPI
	Transformer *transform.PolaroidTransform
	TempDir string
}

type TaskBot struct {
	ChatId int64
	BotError error
	SendBotError bool
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
		TempDir:botConfig.BotTempDir}
}

func (b *Bot) sendTextMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	b.TelegramBot.Send(msg)
}

func (b *Bot) sendPictureMessage(chatId int64, picturePath string) {
	msg := tgbotapi.NewPhotoUpload(chatId, picturePath)
	b.TelegramBot.Send(msg)
}

func (b *Bot) downloadPhoto(chatId int64, photoId, downloadPath string) error {
	log.Printf("Try to download file with %s ID", photoId)
	resp, err := b.TelegramBot.GetFile(tgbotapi.FileConfig{FileID:photoId})
	if err != nil {
		return err
	}
	downloadUrl := telegramRoot + b.TelegramBot.Token + "/" +resp.FilePath
	log.Printf("Handle %s url", downloadUrl)
	raw, err := http.Get(downloadUrl)
	defer raw.Body.Close()
	if err != nil {
		return err
	}
	img, err := jpeg.Decode(raw.Body)
	if err != nil {
		return err
	}
	imaging.Save(img, downloadPath)
	return nil
}

func (b *Bot) BotMainHandler() {
	log.Printf("Authorized on account %s", b.TelegramBot.Self.UserName)
	msgStorage := NewMessagesStorage()

	botTaskChan := make(chan TaskBot, botTaskCap)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := b.TelegramBot.GetUpdatesChan(updateConfig)

	go func() {
		for {
			currentTask := <-botTaskChan
			if currentTask.BotError != nil {
				fullMsg := fmt.Sprintf("Task chat ID - %d Error message - %s", currentTask.ChatId,
					currentTask.BotError.Error())
				log.Println(fullMsg)
				if currentTask.SendBotError {
					log.Println(fmt.Sprintf("Send error message to %d chat ID", currentTask.ChatId))
					b.sendTextMessage(currentTask.ChatId, currentTask.BotError.Error())
				}
			}
		}
	}()

	for update := range updates {
		msgStorage.RemoveExpiredMessages()
		if update.Message == nil {
			continue
		}
		chatId := update.Message.Chat.ID
		if update.Message.Photo != nil {
			srcFileName := utils.GetRandomFileName(b.TempDir)
			dstFileName := utils.GetRandomFileName(b.TempDir)
			msg := msgStorage.GetMessage(chatId)
			go func() {
				task := TaskBot{ChatId:chatId, BotError:nil, SendBotError:false}
				photos := *update.Message.Photo
				photoId := photos[1].FileID
				err := b.downloadPhoto(chatId, photoId, srcFileName)
				if err != nil {
					task.BotError = err
					task.SendBotError = true
				} else {
					err = b.Transformer.CreatePolaroidImage(srcFileName, dstFileName, msg)
					if err != nil {
						task.BotError = err
						task.SendBotError = true
					} else {
						b.sendPictureMessage(chatId, dstFileName)
						err = os.Remove(srcFileName)
						if err != nil {
							task.BotError = err
						}
						err = os.Remove(dstFileName)
						if err != nil {
							task.BotError = err
						}
					}
				}
				botTaskChan <- task
			}()
		}

		switch update.Message.Command() {
		case "start":
			b.sendTextMessage(chatId, botHello)
		case "help":
			b.sendTextMessage(chatId, botHelp)
		default:
			msgStorage.SetMessage(chatId, update.Message.Text)
		}
	}
}
