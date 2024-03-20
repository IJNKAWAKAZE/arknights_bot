package material

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
)

func InlineMaterial(update tgbotapi.Update) error {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "材料-")
	itemList := utils.GetItemsByName(name)
	var inlineQueryResults []interface{}
	for k, v := range itemList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       k,
			Description: "查询" + k,
			ThumbURL:    v,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/material " + k,
			},
		}
		inlineQueryResults = append(inlineQueryResults, queryResult)
	}
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       inlineQueryResults,
		CacheTime:     0,
	}
	bot.Arknights.Send(answerInlineQuery)
	return nil
}
