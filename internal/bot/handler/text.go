package handler

import (
	"context"
	"regexp"

	"github.com/mavihq/persian"
	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
)

var cardNumberRegexPattern = regexp.MustCompile(`^[2569]\d{15}$`)

type text struct {
	db database.Client
}

func NewText(db database.Client) Command {
	return &text{db: db}
}

func (c *text) Endpoint() string {
	return telebot.OnText
}

func (c *text) Handle(ctx telebot.Context) error {
	user := ctx.Get("user").(*entities.User)
	state := user.State
	if state == "setcard" {
		cardNumber := persian.ToEnglishDigits(ctx.Text())
		if !cardNumberRegexPattern.MatchString(cardNumber) {
			return ctx.Reply("شماره کارت نامعتبر است. لطفا شماره کارت را به درستی وارد کنید.")
		}
		user.CardNumber = cardNumber
		user.State = ""
		err := c.db.UpdateUser(context.Background(), user)
		if err != nil {
			return err
		}
		return ctx.Reply("شماره کارت شما با موفقیت ثبت شد.")
	}
	return ctx.Reply("برای راهنمایی کار با بات، از دستور /help استفاده کنید.")
}
