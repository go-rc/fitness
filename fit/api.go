package fit

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func NewFitnessApi(db *mgo.Database, router *mux.Router) {
	entriesCollection := db.C(EntriesCollectionName)
	entriesRepo := NewEntryRepository(entriesCollection)
	entriesController := NewEntriesController(entriesRepo)

	router.HandleFunc("/entries", entriesController.IndexHandler)
}
