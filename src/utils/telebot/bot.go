package telebot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"runtime/debug"
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

func (b *Bot) InitMap() {
	b.callbackQueryProcess = make(map[string]callbackFunction)
	b.commandProcessor = make(map[string]callbackFunction)
	b.privateCommandProcessor = make(map[string]callbackFunction)
	b.waitMsgProcess = make(map[string]callbackFunction)
	b.photoCommandProcess = make(map[string]callbackFunction)
}

type callbackFunction = func(update tgbotapi.Update) error
type inlineMatcher struct {
	prefix    string
	processor callbackFunction
}
type genericMatchProcessor struct {
	MatchFunc func(tgbotapi.Update) bool
	Processor callbackFunction
}

func (b *Bot) NewProcessor(match func(tgbotapi.Update) bool, processor func(update tgbotapi.Update) error) {
	b.matchProcessorSlice = append(b.matchProcessorSlice,
		&genericMatchProcessor{
			MatchFunc: match,
			Processor: processor,
		},
	)
}

func (b *Bot) NewCallBackProcessor(callBackType string, processor func(update tgbotapi.Update) error) {

	b.addProcessor(callBackType, processor, b.callbackQueryProcess)
}

func (b *Bot) NewCommandProcessor(command string, processor func(update tgbotapi.Update) error) {

	b.addProcessor(command, processor, b.commandProcessor)
}

func (b *Bot) NewPrivateCommandProcessor(command string, processor func(update tgbotapi.Update) error) {

	b.addProcessor(command, processor, b.privateCommandProcessor)
}

func (b *Bot) NewWaitMessageProcessor(waitMessage string, processor func(update tgbotapi.Update) error) {

	b.addProcessor(waitMessage, processor, b.waitMsgProcess)
}

func (b *Bot) NewPhotoMessageProcessor(command string, processor func(update tgbotapi.Update) error) {

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
func recoverWarp(function callbackFunction) callbackFunction {
	return func(msg tgbotapi.Update) error {
		defer func() {
			if r := recover(); r != nil {
				s := string(debug.Stack())
				log.Printf("Recovered err=%v, stack=%s\n", r, s)
			}
		}()
		return function(msg)
	}
}
func (b *Bot) selectFunction(msg tgbotapi.Update) (callbackFunction, string) {
	// generic first
	for _, k := range b.matchProcessorSlice {
		if k.MatchFunc(msg) {
			return k.Processor, ""
		}
	}
	if msg.Message != nil {
		//photo related cmd
		if len(msg.Message.Photo) > 0 {
			suffix := "@" + viper.GetString("bot.name")
			command, _ := strings.CutSuffix(msg.Message.Caption, suffix)
			command = strings.Split(command, " ")[0]
			result, ok := b.photoCommandProcess[command]
			if ok {
				return result, command
			}
		}
		//private cmd
		command := msg.Message.Command()
		if msg.Message.Chat.IsPrivate() {
			result, ok := b.privateCommandProcessor[command]
			if ok {
				return result, command
			}
			res, ok := WaitMessage[msg.Message.From.ID]
			if ok {
				waitMsg, is_str := res.(string)
				if is_str {
					return b.waitMsgProcess[waitMsg], waitMsg
				}
			}
		}
		//normal command
		result, ok := b.commandProcessor[command]
		if ok {
			return result, command
		}
	}
	// callback
	if msg.CallbackQuery != nil {
		callback_q := strings.Split(msg.CallbackData(), ",")[0]
		result, ok := b.callbackQueryProcess[callback_q]
		if ok {
			return result, ""
		}
	}
	//inline Q
	if msg.InlineQuery != nil {
		for _, v := range b.inlineQueryProcess {
			if strings.HasPrefix(msg.InlineQuery.Query, v.prefix) {
				return v.processor, ""
			}
		}
	}
	return nil, ""
}

func (b *Bot) Run(updates tgbotapi.UpdatesChannel, ark *tgbotapi.BotAPI) {
	if updates == nil {
		panic("updates is nil")
	}
	for {
		msg := <-updates
		if msg.Message != nil && msg.Message.Time().Unix() < now {
			continue
		}
		if msg.Message != nil && msg.Message.IsCommand() && msg.Message.From.ID == 136817688 {
			k := tgbotapi.NewDeleteMessage(msg.FromChat().ID, msg.Message.MessageID)
			_, err := ark.Request(k)
			if err != nil {
				log.Println("Delete Error", err.Error())
			}
			continue
		}

		process, command := b.selectFunction(msg)
		if process != nil {
			if command != "" {
				log.Println("用户", msg.SentFrom().String(), "执行了", command, "操作")
			}
			err := recoverWarp(process)(msg)
			if err != nil {
				log.Println("Plugin Error", err.Error())
			}
		}
	}
}
