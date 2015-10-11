package main

import (
    "log"
    "github.com/Syfaro/telegram-bot-api"
    // "gopkg.in/mgo.v2"
    "regexp"
    "strings"  
    // "path/filepath"
    "io/ioutil"
    "github.com/valiknet18/AudioSenderTelegramBot/config"
    "github.com/valiknet18/AudioSenderTelegramBot/music"
    // "fmt"
)

func main() {
    session := getSession()
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
        resultRegExp := r.FindStringSubmatch(strings.ToLower(update.Message.Text))

        log.Println(resultRegExp)

        if len(resultRegExp) > 0 {
            switch resultRegExp[2] {
                case "to music": {
                    genre, group, track := music.RandomTrack(session)

                    sendAudioToServer(bot, update, genre, group, track)
                }

                case "genre": {
                    genre, group, track := music.GetRandomTrackByGenre(strings.ToLower(strings.Trim(resultRegExp[3], " ")), session)

                    sendAudioToServer(bot, update, genre, group, track)
                }

                case "group": {
                    genre, group, track := music.GetRandomTrackByGroupName(strings.ToLower(strings.Trim(resultRegExp[3], " ")), session)

                    sendAudioToServer(bot, update, genre, group, track)
                }

                case "track": {
                    genre, group, track := music.GetTrackByTrackNameAndGroupName(strings.ToLower(strings.Trim(resultRegExp[3], " ")), session)

                    sendAudioToServer(bot, update, genre, group, track)
                }
            }
        }   
    }
}

func sendAudioToServer(bot *tgbotapi.BotAPI, update tgbotapi.Update, genre *music.Genre, group *music.Group, track *music.Track) {
    if track == nil {
        waitMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "No one track is not found")
        bot.SendMessage(waitMessage)     
    } else {
        audioBytes, err := ioutil.ReadFile(track.PathToTrack)

        if err != nil {
            panic(err)
        }

        waitMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Wait please, track upload to server")
        bot.SendMessage(waitMessage)

        audio := tgbotapi.FileBytes{Name: strings.Replace(track.NameTrack, " ", "", -1) + ".ogg", Bytes: audioBytes}

        audioConfig := tgbotapi.NewAudioUpload(update.Message.Chat.ID, audio)

        resultMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Result track: " + group.Name + " - " + track.NameTrack)
        bot.SendMessage(resultMessage)

        bot.SendAudio(audioConfig)
    }
}