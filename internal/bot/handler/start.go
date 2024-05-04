package handler

import (
	"strings"

	"github.com/itsamirhn/dongetobede/internal/database"
	"gopkg.in/telebot.v3"
)

type start struct {
	db database.Client
}

func NewStart(db database.Client) Command {
	return &start{db: db}
}

func (c *start) Endpoint() string {
	return "/start"
}

func (c *start) Handle(ctx telebot.Context) error {
	data := ctx.Data()
	if strings.HasPrefix(data, "setcard") {
		return NewSetCard(c.db).Handle(ctx)
	}
	return ctx.Reply("به بات دونگ خوش آمدید!")
}
