package music

import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
    // "fmt"
)

type Group struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Tracks []*Track `bson:"tracks"`
}

func GetRandomTrackByGroupName(groupName string, session *mgo.Session) (*Genre, *Group, *Track) {
	sess := session.DB("audio_sender_telegram").C("genres")

	var genres [] *Genre
	var genre *Genre

	sess.Find(nil).Select(bson.M{"groups": bson.M{"$elemMatch": bson.M{"name": groupName}}}).All(&genres)
	// sess.Find(nil).Select(bson.M{"name": "metal"}).All(&genres)

	if genres == nil {
		return nil, nil, nil
	}

	for _, genreIt := range genres {
		if len(genreIt.Groups) > 0 {
			genre = genreIt
			break
		}
	}

	randomTrack := rand.Intn(len(genre.Groups[0].Tracks))

	return genre, genre.Groups[0], genre.Groups[0].Tracks[randomTrack]
}