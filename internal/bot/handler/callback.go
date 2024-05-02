package handler

import (
	"context"
	"slices"
	"strings"

	"github.com/itsamirhn/dongetobede/internal/database/entities"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v3"
)

type callback struct {
	db database.Client
}

func NewCallback(db database.Client) Command {
	return &callback{db: db}
}

func (c *callback) Endpoint() string {
	return telebot.OnCallback
}

func (c *callback) Handle(ctx telebot.Context) error {
	cb := ctx.Callback()
	prefix := "\fpaydong|"
	data, _ := strings.CutPrefix(cb.Data, prefix)
	dongID, err := primitive.ObjectIDFromHex(data)
	if err != nil {
		return errors.Wrap(err, "invalid dong id")
	}
	dong, err := c.db.GetDongByID(context.Background(), dongID)
	if err != nil {
		return errors.Wrap(err, "failed to get dong by id")
	}
	if dong == nil {
		return errors.New("dong not found")
	}
	user := ctx.Get("user").(*entities.User)
	if dong.IssuerUserID == user.ID {
		return ctx.Respond(&telebot.CallbackResponse{
			Text: "تو که مادر خرجی!",
		})
	}
	if slices.Contains(dong.PaidUserIDs, user.ID) {
		newPaidUserIDs := make([]int64, 0, len(dong.PaidUserIDs))
		for _, id := range dong.PaidUserIDs {
			if id != user.ID {
				newPaidUserIDs = append(newPaidUserIDs, id)
			}
		}
		dong.PaidUserIDs = newPaidUserIDs
		err = ctx.Respond(&telebot.CallbackResponse{
			Text: "دونگ شما پس گرفته شد!",
		})
		if err != nil {
			return errors.Wrap(err, "failed to respond")
		}
	} else {
		dong.PaidUserIDs = append(dong.PaidUserIDs, user.ID)
		err = ctx.Respond(&telebot.CallbackResponse{
			Text: "دونگ شما ثبت شد!",
		})
	}
	err = c.db.UpdateDong(context.Background(), dong)
	if err != nil {
		return errors.Wrap(err, "failed to update dong")
	}
	return ctx.Edit(getDongMarkup(len(dong.PaidUserIDs), dong.TotalPeople, dong.CardNumber, dong.ID))
}
