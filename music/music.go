package music

import (
	"github.com/valiknet18/AudioSenderTelegramBot/Godeps/_workspace/src/gopkg.in/mgo.v2"
    "github.com/valiknet18/AudioSenderTelegramBot/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
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
	trackDataSplited := strings.Split(trackData, "-")

	if len(trackDataSplited) == 2 {
		groupName := strings.Trim(strings.ToLower(trackDataSplited[0]), " ")
		trackName := strings.Trim(strings.ToLower(trackDataSplited[1]), " ")

		group := GetGroup(groupName, session)

		if group.Name == "" {
			return nil, nil, nil
		}

		track := group.GetTrack(trackName)	

		if track.PathToTrack == "" {
			return nil, nil, nil
		}

		return nil, group, track
	}

	return nil, nil, nil
} 