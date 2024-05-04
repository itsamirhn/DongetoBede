package handler

import "gopkg.in/telebot.v3"

type help struct{}

func NewHelp() Command {
	return &help{}
}

func (c *help) Endpoint() string {
	return "/help"
}

func (c *help) Handle(ctx telebot.Context) error {
	markup := &telebot.ReplyMarkup{}
	markup.Inline(
		markup.Row(
			markup.Query("مثال", ""),
		),
	)
	return ctx.Reply(`
برای ثبت خرج در هر گروهی، به گروه خود بروید و دستور زیر را تایپ کنید:
 
<code>@dongetobedebot {مبلغ}</code>

و سپس تعداد نفرات را انتخاب کنید.

اگر میخواهید شماره کارت شما در پیام دونگ نوشته شود، از دستور /setcard استفاده کنید.
`, markup, telebot.ModeHTML)
}
