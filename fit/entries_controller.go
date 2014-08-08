package fit

import (
	"encoding/json"
	"net/http"
)

type EntriesController struct {
	repo *EntryRepository
}

func NewEntriesController(r *EntryRepository) *EntriesController {
	return &EntriesController{repo: r}
}

func (c *EntriesController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	entries, _ := c.repo.Find()
	data, _ := json.Marshal(entries)
	w.Write(data)
}

func (c *EntriesController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var entry Entry
	err := decoder.Decode(&entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.repo.Upsert(&entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
