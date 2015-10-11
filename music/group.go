package music

import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "math/rand"
    // "fmt"
    "github.com/valiknet18/AudioSenderTelegramBot/config"
)

type Group struct {
	session *mgo.Database
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Genre mgo.DBRef `bson:"genre"`
	Tracks []*Track `bson:"tracks"`
}

func (g *Group) SetSession(session *mgo.Database) {
	g.session = session
}

func GetGroup(groupName string, session *mgo.Database) *Group{
	sess := session.C("groups")

	group := &Group{}

	sess.Find(bson.M{"name": groupName}).One(&group)
	group.SetSession(session)

	return group
}

func (g *Group) GetTracks() [] *Track{
	sess := g.session.C("tracks")

	var tracks [] *Track

	sess.Find(bson.M{"group.$id": g.Id}).All(&tracks)

	return tracks
}

func (g *Group) GetTrack(NameTrack string) *Track{
	sess := g.session.C("tracks")

	track := &Track{}

	sess.Find(bson.M{"name_track": NameTrack, "group.$id": g.Id}).One(&track)
	track.SetSession(g.session)

	return track
}

func (g *Group) InsertTrack() *Track{
	sess := g.session.C("tracks")

	configFile := config.ParseConfig()

	track := &Track{session: g.session, NameTrack: "", Group: mgo.DBRef{Collection: "groups", Id: g.Id, Database: configFile.Database}}

	sess.Insert(track)

	return track
}

func GetRandomTrackByGroupName(groupName string, session *mgo.Database) (*Genre, *Group, *Track) {
	group := GetGroup(groupName, session)

	if group == nil {
		return nil, nil, nil
	}

	tracks := group.GetTracks()
	randomTrack := rand.Intn(len(tracks))

	return nil, group, tracks[randomTrack]
}
