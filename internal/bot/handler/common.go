package handler

import (
	"fmt"

	"github.com/mavihq/persian"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v3"
)

func getDongMarkup(paidUsersCount, totalPeople int, cardNumber string, dongID *primitive.ObjectID) *telebot.ReplyMarkup {
	paidUsersCountStr := persian.ToPersianDigitsFromInt(paidUsersCount)
	totalPeopleStr := persian.ToPersianDigitsFromInt(totalPeople)
	btnText := fmt.Sprintf("دنگمو دادم (%s/%s)", paidUsersCountStr, totalPeopleStr)
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, 2)
	if dongID == nil {
		rows = append(rows, markup.Row(markup.Data(btnText, "paydong")))
	} else {
		rows = append(rows, markup.Row(markup.Data(btnText, "paydong", dongID.Hex())))
	}
	markup.Inline(rows...)
	return markup
}
