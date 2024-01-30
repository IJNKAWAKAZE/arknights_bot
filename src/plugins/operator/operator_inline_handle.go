package operator

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func InlineOperator(inlineQuery *tgbotapi.InlineQuery) {
	name := inlineQuery.Query
	operatorList := utils.GetOperatorsByName(name)
	var inlineQueryResults []interface{}
	for _, operator := range operatorList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       operator.Name,
			Description: "查询" + operator.Name,
			ThumbURL:    operator.Painting,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "https://prts.wiki/w/" + operator.Name,
			},
		}
		inlineQueryResults = append(inlineQueryResults, queryResult)
	}
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       inlineQueryResults,
		CacheTime:     0,
	}
	bot.Arknights.Send(answerInlineQuery)
}
