package handler

import (
	"fmt"
	"math"

	"github.com/mavihq/persian"
	"gopkg.in/telebot.v3"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"github.com/itsamirhn/dongetobede/internal/expression"
)

type query struct {
	db        database.Client
	evaluator expression.Evaluator
}

func NewQuery(db database.Client, evaluator expression.Evaluator) Command {
	return &query{
		db:        db,
		evaluator: evaluator,
	}
}

func (c *query) Endpoint() string {
	return telebot.OnQuery
}

func (c *query) getValidLimitedArticle(amount, totalPeople int, cardNumber string) *telebot.ArticleResult {
	perPersonStr := getDongPerPersonToman(amount, totalPeople)
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

func (c *query) getValidUnlimitedArticle(amount int, cardNumber string) *telebot.ArticleResult {
	perPersonStr := getDongPerPersonToman(amount, 0)
	txt := getDongText(amount, 0, cardNumber, nil)
	res := &telebot.ArticleResult{
		Title:       "دنگ نامحدود",
		Description: fmt.Sprintf("نفری %s", perPersonStr),
		Text:        txt,
	}
	res.SetParseMode(telebot.ModeMarkdown)
	res.SetResultID(fmt.Sprintf("%d-%d", amount, 0))
	res.SetReplyMarkup(getDongMarkup(0, 0, cardNumber, nil))
	return res
}

func (c *query) getInvalidArticle() *telebot.ArticleResult {
	res := &telebot.ArticleResult{
		Title:       "مبلغ نامعتبر",
		Description: "مجموع هزینه را به تومان وارد کنید.",
		Text:        ` مبلغ نامعتبر است. لطفا مبلغ را به صورت  عددی یا یک عبارت ریاضی غیر اعشاری وارد کنید.`,
	}
	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(markup.QueryChat("مثال", "56000 + 12000")),
	)
	res.SetReplyMarkup(markup)
	res.SetParseMode(telebot.ModeMarkdown)
	return res
}

func (c *query) getValidArticles(amount int, cardNumber string) []telebot.Result {
	var results []telebot.Result
	results = append(results, c.getValidUnlimitedArticle(amount, cardNumber))
	for i := 2; i <= 15; i++ {
		results = append(results, c.getValidLimitedArticle(amount, i, cardNumber))
	}
	return results
}

func (c *query) Handle(ctx telebot.Context) error {
	q := ctx.Query().Text
	expressionStr := persian.ToEnglishDigits(q)
	amountFloat, err := c.evaluator.Eval(expressionStr)
	amount := int(math.Round(amountFloat))
	if amount <= 0 || err != nil {
		return ctx.Answer(&telebot.QueryResponse{
			Results: []telebot.Result{c.getInvalidArticle()},
		})
	}
	user := ctx.Get("user").(*entities.User)
	return ctx.Answer(&telebot.QueryResponse{
		Results: c.getValidArticles(amount, user.CardNumber),
	})
}
