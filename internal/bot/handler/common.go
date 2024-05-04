package handler

import (
	"fmt"
	"math"
	"strconv"

	"github.com/mavihq/persian"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/database/entities"
)

func getDongMarkup(
	paidUsersCount, totalPeople int,
	_ string,
	dongID *primitive.ObjectID,
) *telebot.ReplyMarkup {
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

func getDongText(amount, totalPeople int, cardNumber string, paidUsers []*entities.User) string {
	perPerson := int(math.Ceil(float64(amount) / float64(totalPeople)))
	perPersonStr := persian.Toman(strconv.Itoa(perPerson))
	txt := fmt.Sprintf("نفری %s", perPersonStr)
	if cardNumber != "" {
		txt += fmt.Sprintf("\nشماره کارت: `%s`", persian.ToEnglishDigits(cardNumber))
	}
	if len(paidUsers) > 0 {
		txt += "\n\nکسایی که دنگشونو دادن:"
		for _, user := range paidUsers {
			identifier := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
			if user.Username != "" {
				identifier = fmt.Sprintf("@%s", user.Username)
			}
			txt += fmt.Sprintf("\n[%s](tg://user?id=%d)", identifier, user.ID)
		}
	}
	return txt
}
