package enemy

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strings"
)

func InlineEnemy(update tgbotapi.Update) (bool, error) {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "敌人-")
	enemyList := utils.GetEnemiesByName(name)
	var inlineQueryResults []interface{}
	for k, v := range enemyList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       k,
			Description: "查询" + k,
			ThumbURL:    v,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/enemy " + k,
			},
		}
		inlineQueryResults = append(inlineQueryResults, queryResult)
	}
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       inlineQueryResults,
		CacheTime:     0,
	}
	_, err := bot.Arknights.Send(answerInlineQuery)
	if err != nil {
		return true, err
	}
	return true, nil
}
