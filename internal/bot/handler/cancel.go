package handler

import (
	"context"

	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
)

type cancel struct {
	db database.Client
}

func NewCancel(db database.Client) Command {
	return &cancel{db: db}
}

func (c *cancel) Endpoint() string {
	return "/cancel"
}

func (c *cancel) Handle(ctx telebot.Context) error {
	user := ctx.Get("user").(*entities.User)
	user.State = ""
	err := c.db.UpdateUser(context.Background(), user)
	if err != nil {
		return err
	}
	return ctx.Reply("عملیات لغو شد.")
}
