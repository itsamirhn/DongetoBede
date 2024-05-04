package handler

import (
	"fmt"
	"math"
	"strconv"

	"github.com/mavihq/persian"
	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
)

type query struct {
	db database.Client
}

func NewQuery(db database.Client) Command {
	return &query{db: db}
}

func (c *query) Endpoint() string {
	return telebot.OnQuery
}

func (c *query) getValidArticle(amount, totalPeople int, cardNumber string) *telebot.ArticleResult {
	perPerson := int(math.Ceil(float64(amount) / float64(totalPeople)))
	perPersonStr := persian.Toman(strconv.Itoa(perPerson))
	totalPeopleStr := persian.ToPersianDigitsFromInt(totalPeople)
	txt := getDongText(amount, totalPeople, cardNumber, nil)
	res := &telebot.ArticleResult{
		Title:       fmt.Sprintf("دنگ %s نفره", totalPeopleStr),
		Description: fmt.Sprintf("نفری %s", perPersonStr),
		Text:        txt,
	}
	res.SetParseMode(telebot.ModeMarkdown)
	res.SetResultID(fmt.Sprintf("%d-%d", amount, totalPeople))
	res.SetReplyMarkup(getDongMarkup(0, totalPeople, cardNumber, nil))
	return res
}

func (c *query) getInvalidArticle() *telebot.ArticleResult {
	res := &telebot.ArticleResult{
		Title:       "مبلغ نامعتبر",
		Description: "مجموع هزینه را به تومان وارد کنید.",
		Text:        ` مبلغ نامعتبر است. لطفا مبلغ را به صورت عددی وارد کنید.`,
	}
	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(markup.QueryChat("مثال", "56000")),
	)
	res.SetReplyMarkup(markup)
	res.SetParseMode(telebot.ModeMarkdown)
	return res
}

func (c *query) getValidArticles(amount int, cardNumber string) []telebot.Result {
	var results []telebot.Result
	for i := 2; i <= 15; i++ {
		results = append(results, c.getValidArticle(amount, i, cardNumber))
	}
	return results
}

func (c *query) Handle(ctx telebot.Context) error {
	q := ctx.Query().Text
	amountStr := persian.ToEnglishDigits(q)
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return ctx.Answer(&telebot.QueryResponse{
			Results: []telebot.Result{c.getInvalidArticle()},
		})
	}
	user := ctx.Get("user").(*entities.User)
	return ctx.Answer(&telebot.QueryResponse{
		Results: c.getValidArticles(amount, user.CardNumber),
	})
}
