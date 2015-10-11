package music

import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
    "strings"
    "fmt"
)

type Track struct {
	session *mgo.Database
	Id bson.ObjectId `bson:"_id"`
	NameTrack string `bson:"name_track"`
	PathToTrack string `bson:"path_to_track"`
	Group mgo.DBRef   `bson:"group"` 
}

func (t *Track) SetSession(session *mgo.Database) {
	t.session = session
}

func (t *Track) GetGroup() *Group{
	sess := t.session.C("groups")

	group := &Group{session: t.session}

	sess.Find(bson.M{"_id": t.Group.Id}).One(&group)

	return group
}

func RandomTrack(session *mgo.Database) (*Genre, *Group, *Track) {
	sess := session.C("tracks")

	track := &Track{}

	count, _ := sess.Count()

	sess.Find(bson.M{}).Limit(-1).Skip(rand.Intn(count)).One(&track)
	track.SetSession(session)

	if track == nil {
		return nil, nil, nil
	}

	fmt.Println(track)

	group := track.GetGroup()

	return nil, group, track
}

func GetTrackByTrackNameAndGroupName(trackData string, session *mgo.Database) (*Genre, *Group, *Track) {
	sess := session.C("genres")

	trackDataSplited := strings.Split(trackData, "-")

	if len(trackDataSplited) == 2 {
		groupName := strings.Trim(strings.ToLower(trackDataSplited[0]), " ")
		trackName := strings.Trim(strings.ToLower(trackDataSplited[1]), " ")

		var genres [] *Genre
		// var genre *Genre
		// var group *Group

		sess.Find(nil).Select(bson.M{"groups": bson.M{"$elemMatch": bson.M{"name": groupName, "tracks": bson.M{"$elemMatch": bson.M{"name_track": trackName}}}}}).All(&genres)

		// for _, genreIt := range genres {
		// 	if len(genreIt.Groups) > 0 {
		// 		genre = genreIt
		// 		break
		// 	}
		// }

		// if genre == nil {
		// 	return nil, nil, nil
		// }

		// for _, groupIt := range genre.Groups {
		// 	if len(groupIt.Tracks) > 0 {
		// 		group = groupIt
		// 		break;
		// 	}
		// }

		return nil, nil, nil	
	}

	return nil, nil, nil
} 