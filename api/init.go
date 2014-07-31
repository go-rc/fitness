package api

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

const EntriesCollectionName string = "entries"

func NewFitnessApi(db *mgo.Database, router *mux.Router) {
	entriesCollection := db.C(EntriesCollectionName)
	entriesRepo := NewEntryRepository(entriesCollection)
	entriesController := NewEntriesController(entriesRepo)

	router.HandleFunc("/entries", entriesController.IndexHandler)
}
