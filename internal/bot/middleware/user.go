package middleware

import (
	"context"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

func NewUserRetriever(db database.Client) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			sender := ctx.Sender()
			if sender != nil {
				user, err := db.GetUserByID(context.Background(), sender.ID)
				if errors.Is(err, database.ErrNotFound) {
					user = &entities.User{
						ID:        sender.ID,
						Username:  sender.Username,
						FirstName: sender.FirstName,
						LastName:  sender.LastName,
					}
					_, err := db.AddUser(context.Background(), user)
					if err != nil {
						return err
					}
				}
				ctx.Set("user", user)
			}
			return next(ctx)
		}
	}
}
