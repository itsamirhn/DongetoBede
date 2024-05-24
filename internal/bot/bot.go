package bot

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/bot/handler"
	"github.com/itsamirhn/dongetobede/internal/bot/middleware"
	"github.com/itsamirhn/dongetobede/internal/config"
	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/expression"
	"github.com/itsamirhn/dongetobede/pkg"
)

type Bot struct {
	*telebot.Bot
}

func NewBot(db database.Client, evaluator expression.Evaluator,
	token string, endpoint string, listenPort string,
) (*Bot, error) {
	settings := telebot.Settings{
		Token:   token,
		Poller:  getTelebotPoller(endpoint, listenPort),
		Verbose: config.GlobalConfig.Verbose,
	}
	tgBot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot := &Bot{Bot: tgBot}
	bot.Use(middleware.NewUserRetriever(db))
	bot.registerCommands([]handler.Command{
		handler.NewStart(db),
		handler.NewSetCard(db),
		handler.NewQuery(db, evaluator),
		handler.NewInline(db),
		handler.NewCallback(db),
		handler.NewText(db),
		handler.NewHelp(),
		handler.NewCancel(db),
	})

	return bot, nil
}

func (b *Bot) registerCommands(commands []handler.Command) {
	for _, h := range commands {
		b.Handle(h.Endpoint(), h.Handle)
	}
}

func getTelebotPoller(endpoint, listenPort string) telebot.Poller {
	if config.GlobalConfig.DebugMode {
		return &telebot.LongPoller{Timeout: 10 * time.Second}
	}
	return &telebot.Webhook{
		Endpoint:    &telebot.WebhookEndpoint{PublicURL: endpoint},
		Listen:      fmt.Sprintf(":%v", listenPort),
		SecretToken: pkg.RandString(10),
	}
}
