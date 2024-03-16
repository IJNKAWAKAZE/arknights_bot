package telebot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

var WaitMessage = make(map[int64]interface{})

type Bot struct {
	matchProcessorSlice     []*genericMatchProcessor
	commandProcessor        map[string]callbackFunction
	privateCommandProcessor map[string]callbackFunction
	photoCommandProcess     map[string]callbackFunction
	inlineQueryProcess      []*inlineMatcher
	callbackQueryProcess    map[string]callbackFunction
	waitMsgProcess          map[string]callbackFunction
}

var now = time.Now().Unix()

type callbackFunction = func(update tgbotapi.Update) (bool, error) // isBreak,error
type inlineMatcher struct {
	prefix    string
	processor callbackFunction
}
type genericMatchProcessor struct {
	MatchFunc func(tgbotapi.Update) bool
	Processor callbackFunction
}

func (b *Bot) NewProcessor(match func(tgbotapi.Update) bool, processor func(update tgbotapi.Update) (bool, error)) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&genericMatchProcessor{
			MatchFunc: match,
			Processor: processor,
		},
	)
}

func (b *Bot) NewCallBackProcessor(callBackType string, processor func(update tgbotapi.Update) (bool, error)) {

	b.addProcessor(callBackType, processor, b.callbackQueryProcess)
}

func (b *Bot) NewCommandProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {

	b.addProcessor(command, processor, b.commandProcessor)
}

func (b *Bot) NewPrivateCommandProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {

	b.addProcessor(command, processor, b.privateCommandProcessor)
}

func (b *Bot) NewWaitMessageProcessor(waitMessage string, processor func(update tgbotapi.Update) (bool, error)) {

	b.addProcessor(waitMessage, processor, b.waitMsgProcess)
}

func (b *Bot) NewPhotoMessageProcessor(command string, processor func(update tgbotapi.Update) (bool, error)) {

	b.addProcessor(command, processor, b.photoCommandProcess)
}

func (b *Bot) NewInlineQueryProcessor(command string, processor callbackFunction) {
	b.inlineQueryProcess = append(b.inlineQueryProcess, &inlineMatcher{
		prefix:    command,
		processor: processor,
	})
}
func (b *Bot) addProcessor(command string, processor callbackFunction, funcMap map[string]callbackFunction) {
	_, ok := funcMap[command]
	if ok {
		log.Printf("command %s is already added overriding \n", command)
	}
	funcMap[command] = processor
}
func (b *Bot) selectFunction(msg tgbotapi.Update) callbackFunction {
	if msg.Message != nil {
		// wait msg
		if msg.Message.Chat.IsPrivate() {
			res, ok := WaitMessage[msg.Message.From.ID]
			if ok {
				waitMsg, is_str := res.(string)
				if is_str {
					return b.waitMsgProcess[waitMsg]
				}
			}
		}
		//photo related cmd
		if len(msg.Message.Photo) > 0 {
			result, ok := b.photoCommandProcess[msg.Message.Caption]
			if ok {
				return result
			}
		}
		//private cmd
		command := msg.Message.Command()
		if msg.Message.Chat.IsPrivate() {
			result, ok := b.privateCommandProcessor[command]
			if ok {
				return result
			}
		}
		//normal command
		result, ok := b.commandProcessor[command]
		if ok {
			return result
		}
	}
	// callback
	if msg.CallbackQuery != nil {
		callback_q := strings.Split(msg.CallbackData(), ",")[0]
		result, ok := b.callbackQueryProcess[callback_q]
		if ok {
			return result
		}
	}
	//inline Q
	if msg.InlineQuery != nil {
		for _, v := range b.inlineQueryProcess {
			if strings.HasPrefix(msg.InlineQuery.Query, v.prefix) {
				return v.processor
			}
		}
	}
	return nil
}

func (b *Bot) Run(updates tgbotapi.UpdatesChannel) {
	if updates == nil {
		panic("updates is nil")
	}
	for {
		msg := <-updates
		if msg.Message != nil && msg.Message.Time().Unix() < now {
			continue
		}
		process := b.selectFunction(msg)
		isBreak, err := process(msg)
		if err != nil {
			log.Println("Plugin Error", err.Error())
		}
		if isBreak {
			break
		}
	}
}
