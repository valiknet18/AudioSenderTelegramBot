package music

import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
    "strings"
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

func GetTrackByTrackNameAndGroupName(trackData string, session *mgo.Session) (*Genre, *Group, *Track) {
	sess := session.DB("audio_sender_telegram").C("genres")

	trackDataSplited := strings.Split(trackData, "-")

	if len(trackDataSplited) == 2 {
		groupName := strings.Trim(strings.ToLower(trackDataSplited[0]), " ")
		trackName := strings.Trim(strings.ToLower(trackDataSplited[1]), " ")

		var genres [] *Genre
		var genre *Genre
		var group *Group

		sess.Find(nil).Select(bson.M{"groups": bson.M{"$elemMatch": bson.M{"name": groupName, "tracks": bson.M{"$elemMatch": bson.M{"name_track": trackName}}}}}).All(&genres)

		for _, genreIt := range genres {
			if len(genreIt.Groups) > 0 {
				genre = genreIt
				break
			}
		}

		if genre == nil {
			return nil, nil, nil
		}

		for _, groupIt := range genre.Groups {
			if len(groupIt.Tracks) > 0 {
				group = groupIt
				break;
			}
		}

		return genre, group, group.Tracks[0]	
	}

	return nil, nil, nil
} 