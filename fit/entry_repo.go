package fit

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const EntriesCollectionName string = "entries"

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

func (r *EntryRepository) Upsert(e *Entry) error {
	_, err := r.collection.Upsert(bson.M{"date": e.Date}, e)

	if err != nil {
		return err
	}

	return nil
}
