package api

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EntryRepository struct {
	collection *mgo.Collection
}

func NewEntryRepository(c *mgo.Collection) *EntryRepository {
	return &EntryRepository{collection: c}
}

func (r *EntryRepository) Find() (*[]Entry, error) {
	var results []Entry
	err := r.collection.Find(bson.M{}).All(&results)

	if err != nil {
		return nil, err
	}

	return &results, nil
}
