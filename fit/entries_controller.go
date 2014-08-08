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
