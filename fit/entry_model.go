package fit

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entry struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Date             time.Time     `bson:"date,omitempty" json:"date,omitifempty"`
	Weight           float32       `bson:"weight,omitempty" json:"weight,omitifempty"`
	CaloriesConsumed float32       `bson:"calories_consumed,omitempty" json:"calories_consumed,omitifempty"`
	CaloriesBurned   float32       `bson:"calories_burned,omitempty" json:"calories_burned,omitifempty"`
}
