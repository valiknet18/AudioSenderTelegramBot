package music

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"math/rand"
	"github.com/valiknet18/AudioSenderTelegramBot/config"
	"fmt"
)

type Genre struct {
	session *mgo.Database
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
}

func (g *Genre) SetSession(session *mgo.Database) {
	g.session = session
}

func (g *Genre) GetGroups() []*Group{
	sess := g.session.C("groups")

	var groups [] *Group

	sess.Find(bson.M{"genre.$id": g.Id}).All(&groups);

	return groups
}

func (g *Genre) GetGroup(nameGroup string) *Group{
	sess := g.session.C("groups")

	var group *Group

	sess.Find(bson.M{"name": nameGroup}).Select(bson.M{"genre": bson.M{"$id": g.Id}}).One(&group)
	group.SetSession(g.session)

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

	genre := &Genre{}

	sess.Find(bson.M{"name": genreName}).One(&genre)
	genre.SetSession(session)

	if genre == nil {
		return nil, nil, nil
	}

	groups := genre.GetGroups()

	if len(groups) <= 0 {
		return nil, nil, nil
	}

	randomGroup := rand.Intn(len(groups))

	fmt.Println(groups[randomGroup])

	groups[randomGroup].SetSession(session)
	tracks := groups[randomGroup].GetTracks();
	
	fmt.Println(tracks)

	if len(tracks) <= 0 {
		return nil, nil, nil
	}

	randomTrack := rand.Intn(len(tracks))

	return genre, groups[randomGroup], tracks[randomTrack]
}