package main

import (
    "log"
    "github.com/Syfaro/telegram-bot-api"
    // "gopkg.in/mgo.v2"
    "regexp"
    "strings"  
    // "path/filepath"
    "io/ioutil"
    "valiknet/IntertechTelegramBot/config"
)

func main() {
    // _ := getSession()
    config := config.ParseConfig()

    bot, err := tgbotapi.NewBotAPI(config.BotApi)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    err = bot.UpdatesChan(u)
    if err != nil {
        log.Panic(err)
    }

    for update := range bot.Updates {
        r, _ := regexp.Compile("(^i want to listen (to music|group|genre|track)([a-zA-Zа-яА-Я0-9- ]*)$)")
        resultRegExp := r.FindAllString(strings.ToLower(update.Message.Text), -1)

        audioBytes, err := ioutil.ReadFile("static/Denis_Shatskikh_-_Moim_Druzyam.ogg")

        if err != nil {
            panic(err)
        }

        waitMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Wait please, track upload to server")
        bot.SendMessage(waitMessage)

        audio := tgbotapi.FileBytes{Name: "Moim_Druziam.ogg", Bytes: audioBytes}

        audioConfig := tgbotapi.NewAudioUpload(update.Message.Chat.ID, audio)

        bot.SendAudio(audioConfig)
        log.Println(resultRegExp)
    }
}