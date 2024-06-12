package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	tele "gopkg.in/telebot.v3"
)

const webhookURL = "https://api.example.com/tgbot"

type BotManager struct {
	bots     map[string]*tele.Bot
	botsLock sync.RWMutex
}

func NewBotManager() *BotManager {
	return &BotManager{
		bots: make(map[string]*tele.Bot),
	}
}

func (m *BotManager) InitBot(token string) (*tele.Bot, error) {
	m.botsLock.RLock()
	instance, exists := m.bots[token]
	m.botsLock.RUnlock()

	if exists {
		return instance, nil
	}

	pref := tele.Settings{
		Verbose:     true,
		Token:       token,
		Synchronous: true,
		Offline:     false,
		Poller: &tele.Webhook{
			DropUpdates: true,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: fmt.Sprintf("%s/%s", webhookURL, token),
			},
		},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Errorf("Failed to create bot: %v", err)
		return nil, err
	}

	registerBotHandles(bot)

	m.botsLock.Lock()
	m.bots[token] = bot
	m.botsLock.Unlock()

	return bot, nil
}

func registerBotHandles(bot *tele.Bot) {
	bot.Handle("/start", func(c tele.Context) error {
		body := fmt.Sprintf(`Welcome!`) + "\n\n"

		buttons := &tele.ReplyMarkup{}
		buttons.Inline(
			buttons.Row(
				buttons.Data(fmt.Sprintf(`🆕 Create a bot`), "/newbot", "111", "2222"),
				buttons.Data(fmt.Sprintf(`🤖 Manage bots`), "/mybots"),
			),
		)
		return c.Send(body, tele.ModeHTML, tele.NoPreview, buttons)
	})

	bot.Handle("/newbot", OnTest)
}

func OnTest(c tele.Context) error {
	payload := c.Args()
	return c.Send(fmt.Sprintf("test: %v", payload))
}

func TGBot(tgbot *echo.Group) {
	log.Debugf("tgbot router ...")

	tgbot.GET("/:token", func(c echo.Context) error {
		token := c.Param("token")
		bot, err := NewBotManager().InitBot(token)
		if err != nil {
			log.Fatal(err)
			return c.String(500, "Internal Server Error")
		}

		if err := bot.SetWebhook(bot.Poller.(*tele.Webhook)); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		res, err := bot.Webhook()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(200, map[string]interface{}{
			"bot":     bot.Me,
			"webhook": res,
		})
	})

	tgbot.POST("/:token", func(c echo.Context) error {
		update := tele.Update{}
		if err := c.Bind(&update); err != nil {
			return c.String(400, "Bad Request")
		}

		token := c.Param("token")
		bot, err := NewBotManager().InitBot(token)
		if err != nil {
			log.Fatal(err)
			return c.String(500, "Internal Server Error")
		}

		bot.ProcessUpdate(update)

		return c.JSON(200, bot.Response())
	})
}
