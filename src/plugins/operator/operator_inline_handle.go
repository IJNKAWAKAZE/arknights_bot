package operator

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
)

func InlineOperator(update tgbotapi.Update) error {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "干员-")
	operatorList := utils.GetOperatorsByName(name)
	var inlineQueryResults []interface{}
	for _, operator := range operatorList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       operator.Name,
			Description: "查询" + operator.Name,
			ThumbURL:    operator.ThumbURL,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/operator " + operator.Name,
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
