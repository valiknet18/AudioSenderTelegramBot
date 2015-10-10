package music

import (
	"gopkg.in/mgo.v2/bson"
	"valiknet"
)

type Genre struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Groups []*Group `bson:"groups"`
}

func GetRandomGroupByGenre(genre string, session ) {

}