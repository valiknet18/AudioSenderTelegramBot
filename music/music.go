package music

import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
)

type Track struct {
	Id bson.ObjectId `bson:"_id"`
	NameTrack string `bson:"name_track"`
	PathToTrack string `bson:"path_to_track"` 
}

func RandomTrack(session *mgo.Session) (*Genre, *Group, *Track) {
	sess := session.DB("audio_sender_telegram").C("genres")

	var genre *Genre

	count, _ := sess.Count()

	sess.Find(bson.M{}).Limit(-1).Skip(rand.Intn(count)).One(&genre)

	if genre == nil {
		return nil, nil, nil
	}

	randomGroup := rand.Intn(len(genre.Groups)) 
	randomTrack := rand.Intn(len(genre.Groups[randomGroup].Tracks))

	return genre, genre.Groups[randomGroup], genre.Groups[randomGroup].Tracks[randomTrack]
}