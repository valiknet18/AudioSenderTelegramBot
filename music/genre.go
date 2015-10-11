package music

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"math/rand"
)

type Genre struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Groups []*Group `bson:"groups"`
}

func GetRandomTrackByGenre(genreName string, session *mgo.Session) (*Genre, *Group, *Track) {
	sess := session.DB("audio_sender_telegram").C("genres")

	var genre *Genre

	sess.Find(bson.M{"name": genreName}).One(&genre)

	if genre == nil {
		return nil, nil, nil
	}

	randomGroup := rand.Intn(len(genre.Groups)) 
	randomTrack := rand.Intn(len(genre.Groups[randomGroup].Tracks))

	return genre, genre.Groups[randomGroup], genre.Groups[randomGroup].Tracks[randomTrack]
}