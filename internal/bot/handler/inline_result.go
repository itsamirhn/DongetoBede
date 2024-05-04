package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/itsamirhn/dongetobede/internal/database"
	"github.com/itsamirhn/dongetobede/internal/database/entities"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type inline struct {
	db database.Client
}

func NewInline(db database.Client) Command {
	return &inline{db: db}
}

func (c *inline) Endpoint() string {
	return telebot.OnInlineResult
}

func (c *inline) Handle(ctx telebot.Context) error {
	id := ctx.InlineResult().ResultID
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return errors.New("invalid result id")
	}
	amount, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.Wrap(err, "invalid amount")
	}
	totalPeople, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.Wrap(err, "invalid total people")
	}

	user := ctx.Get("user").(*entities.User)
	dong := &entities.Dong{
		Amount:       amount,
		IssuerUserID: user.ID,
		CardNumber:   user.CardNumber,
		TotalPeople:  totalPeople,
		PaidUserIDs:  []int64{user.ID},
		MessageID:    ctx.InlineResult().MessageID,
	}
	dongID, err := c.db.AddDong(context.Background(), dong)
	if err != nil {
		return err
	}
	dong.ID = dongID
	return ctx.Edit(getDongText(dong.Amount, dong.TotalPeople, dong.CardNumber, []*entities.User{user}), telebot.ModeMarkdown, getDongMarkup(len(dong.PaidUserIDs), dong.TotalPeople, dong.CardNumber, dong.ID))
}
