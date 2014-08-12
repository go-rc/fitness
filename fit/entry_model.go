package fit

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entry struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Date             time.Time     `bson:"date,omitempty" json:"date,omitifempty"`
	Weight           float64       `bson:"weight,omitempty" json:"weight,omitifempty"`
	CalorieGoal      float64       `bson:"calorie_goal,omitempty" json:"calorie_goal,omitifempty"`
	CaloriesConsumed float64       `bson:"calories_consumed,omitempty" json:"calories_consumed,omitifempty"`
	CaloriesBurned   float64       `bson:"calories_burned,omitempty" json:"calories_burned,omitifempty"`
	NetCalories      float64       `bson:net_calories,omitempty" json:"net_calories,omitifempty"`
}
