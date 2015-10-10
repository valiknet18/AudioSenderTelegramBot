package music

import (
	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	Id bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	Tracks []*music `bson:"tracks"`
}