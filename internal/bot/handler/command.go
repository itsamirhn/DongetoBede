package handler

import "gopkg.in/telebot.v3"

type Command interface {
	Endpoint() string
	Handle(ctx telebot.Context) error
}
