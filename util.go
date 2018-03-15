package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TimeToMongoTime(t time.Time) bson.MongoTimestamp {
	return bson.MongoTimestamp(t.UTC().Unix() << 32)
}
