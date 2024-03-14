package telebot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

var WaitMessage = make(map[int64]interface{})

type Bot struct {
	matchProcessorSlice []*matchProcessor
}

var now = time.Now().Unix()

type matchProcessor struct {
	MatchFunc func(tgbotapi.Update) bool
	Processor func(update tgbotapi.Update) (bool, error) // isBreak,error
}

func (b *Bot) NewProcessor(match func(tgbotapi.Update) bool, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: match,
			Processor: processor,
		},
	)
}

func (b *Bot) NewCallBackProcessor(callBackType string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.CallbackQuery != nil && strings.Split(update.CallbackData(), ",")[0] == callBackType
			},
			Processor: processor,
		},
	)
}

func (b *Bot) NewCommandProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == command
			},
			Processor: processor,
		},
	)
}

func (b *Bot) NewPrivateCommandProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == command && update.Message.Chat.IsPrivate()
			},
			Processor: processor,
		},
	)
}

func (b *Bot) NewWaitMessageProcessor(waitMessage string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.Message != nil && update.Message.Chat.IsPrivate() && WaitMessage[update.Message.From.ID] == waitMessage
			},
			Processor: processor,
		},
	)
}

func (b *Bot) NewPhotoMessageProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.Message != nil && len(update.Message.Photo) > 0 && update.Message.Caption == command
			},
			Processor: processor,
		},
	)
}

func (b *Bot) NewInlineQueryProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&matchProcessor{
			MatchFunc: func(update tgbotapi.Update) bool {
				return update.InlineQuery != nil && strings.HasPrefix(update.InlineQuery.Query, command)
			},
			Processor: processor,
		},
	)
}

func (b *Bot) Run(updates tgbotapi.UpdatesChannel) {
	if updates == nil {
		panic("updates is nil")
	}
	for {
		select {
		case msg := <-updates:
			//log.Printf("Bot.Run%#v\n", msg)
			if msg.Message != nil && msg.Message.Time().Unix() < now {
				continue
			}
			for i := range b.matchProcessorSlice {
				if b.matchProcessorSlice[i].MatchFunc(msg) {
					log.Println("Bot.Match", i)
					isBreak, err := b.matchProcessorSlice[i].Processor(msg)
					if err != nil {
						log.Println("Bot.Match Error", err.Error())
					}
					if isBreak {
						break
					}
				}
			}
		}
	}
}
