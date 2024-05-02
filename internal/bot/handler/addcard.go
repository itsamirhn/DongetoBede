package handler

import (
	"context"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"gopkg.in/telebot.v3"
)

type addCard struct {
	db database.Client
}

func NewAddCard(db database.Client) Command {
	return &addCard{db: db}
}

func (c *addCard) Endpoint() string {
	return "/addcard"
}

func (c *addCard) Handle(ctx telebot.Context) error {
	user := ctx.Get("user").(*entities.User)
	state := "addcard"
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
