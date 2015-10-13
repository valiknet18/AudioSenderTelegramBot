package main

import (
	"github.com/valiknet18/AudioSenderTelegramBot/Godeps/_workspace/src/github.com/Syfaro/telegram-bot-api"
	"github.com/valiknet18/AudioSenderTelegramBot/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"log"
	"regexp"
	"strings"
	// "path/filepath"
	"encoding/json"
	"fmt"
	"github.com/valiknet18/AudioSenderTelegramBot/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/valiknet18/AudioSenderTelegramBot/config"
	"github.com/valiknet18/AudioSenderTelegramBot/music"
	"io/ioutil"
	"net/http"
)

var session *mgo.Database

func main() {
	session = getDatabase()

	go controlBot(session)

	r := mux.NewRouter()
	r.HandleFunc("/", IndexAction)
	r.HandleFunc("/genre/create", GetGenres)
	r.HandleFunc("/group", GetGroups)
	r.HandleFunc("/group/{id}", GetTracks)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)

	fmt.Println("Server was started on port :8000")
	http.ListenAndServe(":8000", nil)
}

func IndexAction(rw http.ResponseWriter, req *http.Request) {
	page, err := ioutil.ReadFile("static/index.html")

	if err != nil {
		log.Fatal("Static page not loaded")
	}

	fmt.Fprintf(rw, string(page))
}

func GetGenres(res http.ResponseWriter, req *http.Request) {
	sess := session.C("genres")

	var genres []*music.Genre

	sess.Find(nil).All(&genres)

	jsonGenres, err := json.Marshal(genres)

	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-type", "application/json")
	res.Write(jsonGenres)
}

func GetGroups(res http.ResponseWriter, req *http.Request) {

}

func GetTracks(res http.ResponseWriter, req *http.Request) {

}

func controlBot(session *mgo.Database) {
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
			case "to music":
				{
					genre, group, track := music.RandomTrack(session)

					sendAudioToServer(bot, update, genre, group, track)
				}

			case "genre":
				{
					genre, group, track := music.GetRandomTrackByGenre(strings.ToLower(strings.Trim(resultRegExp[3], " ")), session)

					sendAudioToServer(bot, update, genre, group, track)
				}

			case "group":
				{
					genre, group, track := music.GetRandomTrackByGroupName(strings.ToLower(strings.Trim(resultRegExp[3], " ")), session)

					sendAudioToServer(bot, update, genre, group, track)
				}

			case "track":
				{
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
		audioBytes, err := ioutil.ReadFile("static/music/" + track.PathToTrack)

		if err != nil {
			panic(err)
		}

		waitMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Wait please, track upload to server")
		bot.SendMessage(waitMessage)

		audio := tgbotapi.FileBytes{Name: strings.Replace(track.NameTrack, " ", "", -1) + ".ogg", Bytes: audioBytes}

		audioConfig := tgbotapi.NewAudioUpload(update.Message.Chat.ID, audio)

		resultMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Result track: "+group.Name+" - "+track.NameTrack)
		bot.SendMessage(resultMessage)

		bot.SendAudio(audioConfig)
	}
}
