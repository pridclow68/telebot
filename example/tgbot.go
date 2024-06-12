package main

import (
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	tele "gopkg.in/telebot.v3"
)

const webhookURL = "https://local.test.example.com/tgbot"

type BotInstance struct {
	bot  *tele.Bot
	once sync.Once
}

var (
	bots     = make(map[string]*BotInstance)
	botsLock = sync.RWMutex{}
)

func initBot(token string) (*BotInstance, error) {
	botsLock.RLock()
	instance, exists := bots[token]
	botsLock.RUnlock()

	if exists {
		return instance, nil
	}

	// 创建 Telebot 设置
	pref := tele.Settings{
		Verbose:     true,
		Token:       token,
		Synchronous: true,
		Offline:     true,
		Poller: &tele.Webhook{
			DropUpdates: true,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: fmt.Sprintf("%s/%s", webhookURL, token),
			},
		},
	}
	// 创建 Bot 实例
	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	instance = &BotInstance{
		bot: bot,
	}

	botsLock.Lock()
	bots[token] = instance
	botsLock.Unlock()

	return instance, nil
}

func TGBot(tgbot *echo.Group) {
	log.Debugf("tgbot router ...")

	tgbot.GET("/:token", func(c echo.Context) error {
		token := c.Param("token")
		instance, err := initBot(token)
		if err != nil {
			log.Fatal(err)
			return c.String(500, "Internal Server Error")
		}

		stop := make(chan struct{})
		go instance.bot.Poller.Poll(instance.bot, instance.bot.Updates, stop)
		stop <- struct{}{}

		return c.String(200, "Hello, Telegram Bot!")
	})

	tgbot.POST("/:token", func(c echo.Context) error {
		// 先验证数据格式
		update := tele.Update{}
		if err := c.Bind(&update); err != nil {
			return c.String(400, "Bad Request")
		}

		// 设置 Token
		token := c.Param("token")

		instance, err := initBot(token)
		if err != nil {
			log.Fatal(err)
			return c.String(500, "Internal Server Error")
		}

		instance.once.Do(func() {
			log.Infof("tgbot %s start ...", token)
			instance.bot.Handle("/hello", func(c tele.Context) error {
				// resp, err := instance.bot.Send(c.Chat(), "Hello!")
				// if err != nil {
				// 	log.Errorf("tgbot %s send error: %v", token, err)
				// } else {
				// 	log.Infof("tgbot %s send: %v", token, resp)
				// }
				// return c.Send("Hello!", tele.UseWebhook)
				return c.Delete()
			})
		})
		instance.bot.ProcessUpdate(update)

		return c.JSON(200, instance.bot.Response())
	})
}
