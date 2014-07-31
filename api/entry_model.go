package api

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entry struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Date             time.Time     `bson:"date,omitempty"`
	Weight           float32       `bson:"weight,omitempty"`
	CaloriesConsumed float32       `bson:"calories_consumed,omitempty"`
	CaloriesBurned   float32       `bson:"calories_burned,omitempty"`
}
