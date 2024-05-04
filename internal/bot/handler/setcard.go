package handler

import (
	"context"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"gopkg.in/telebot.v3"
)

type setCard struct {
	db database.Client
}

func NewSetCard(db database.Client) Command {
	return &setCard{db: db}
}

func (c *setCard) Endpoint() string {
	return "/setcard"
}

func (c *setCard) Handle(ctx telebot.Context) error {
	user := ctx.Get("user").(*entities.User)
	state := "setcard"
	if ctx.Data() != "" {
		state = ctx.Data()
	}
	user.State = state
	err := c.db.UpdateUser(context.Background(), user)
	if err != nil {
		return err
	}
	return ctx.Reply("برای ست شدن پیش‌فرض شماره کارت، لطفا شماره کارت خود را ارسال کنید.")
}
