package music

import (
    "gopkg.in/mgo.v2/bson"
)

type Music struct {
	Id bson.ObjectId `bson:"_id"`
	NameTrack string `bson:"name_track"`
	PathToTrack string `bson:"path_to_track"` 
}