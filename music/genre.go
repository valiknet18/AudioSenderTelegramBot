package music

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"math/rand"
	"github.com/valiknet18/AudioSenderTelegramBot/config"
)

type Genre struct {
	session *mgo.Database
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
}

func (g *Genre) GetGroups() []*Group{
	sess := g.session.C("groups")

	var groups [] *Group

	sess.Find(nil).Select(bson.M{"ref": bson.M{"$id": g.Id}}).All(groups);

	return groups
}

func (g *Genre) GetGroup(nameGroup string) *Group{
	sess := g.session.C("groups")

	var group *Group

	sess.Find(bson.M{"name": nameGroup}).Select(bson.M{"ref": bson.M{"$id": g.Id}}).One(&group)

	return group
}

func (g *Genre) InsertGroup() *Group{
	sess := g.session.C("groups")

	configFile := config.ParseConfig()

	group := &Group{Name: "", Genre: mgo.DBRef{Collection: "genres", Id: g.Id, Database: configFile.Database}}

	sess.Insert(group)

	return group
}

func GetRandomTrackByGenre(genreName string, session *mgo.Database) (*Genre, *Group, *Track) {
	sess := session.C("genres")

	genre := &Genre{session: session}

	sess.Find(bson.M{"name": genreName}).One(&genre)

	if genre == nil {
		return nil, nil, nil
	}

	groups := genre.GetGroups()
	randomGroup := rand.Intn(len(groups))

	tracks := groups[randomGroup].GetTracks();
	randomTrack := rand.Intn(len(tracks))

	return genre, groups[randomGroup], tracks[randomTrack]
}