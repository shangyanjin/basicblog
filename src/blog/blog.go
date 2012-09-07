// Package blog contains blog and blog entry types and functions for loading,
// saving and using blogs
package blog

import (
	"encoding/json"
	"os"
	"time"
)

// Type alias for sort.Interface
type BlogEntries []BlogEntry
type ByID struct{ BlogEntries }
type ByDate struct{ BlogEntries }

// Blog represents a blog with all of the entries belonging to it.
type Blog struct {
	Entries BlogEntries
}

// BlogEntry represents a single blog entry.
type BlogEntry struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
}

// NewFromFile attempts to load a blog from the given file fn.
func NewFromFile(fn string) (b *Blog, err error) {
	f, err := os.Open(fn)
	if err != nil {
		return
	}

	defer f.Close()

	var tmp Blog

	d := json.NewDecoder(f)
	if err = d.Decode(&tmp); err != nil {
		return
	}

	return &tmp, nil
}

// Save attempts to write this blog to the given file fn.
func (b *Blog) Save(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(b)
}

// AddEntry adds a single entry to the front of this blog.
func (b *Blog) AddEntry(e *BlogEntry) {
	var highestId int = 0

	for _, e := range b.Entries {
		if e.ID > highestId {
			highestId = e.ID
		}
	}

	nextId := highestId + 1
	e.ID = nextId

	b.Entries = append(b.Entries, *e)
}

// Implements sort.Interface.Len()
func (e BlogEntries) Len() int {
	return len(e)
}

// Implements sort.Interface.Swap()
func (e BlogEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Implements sort.Interface.Less(), sorting the posts by descending date
func (e ByDate) Less(i, j int) bool {
	return e.BlogEntries[i].Date.After(e.BlogEntries[j].Date)
}

// Implements sort.Interface.Less(), sorting the posts by ascending IDs
func (e ByID) Less(i, j int) bool {
	return e.BlogEntries[i].ID < e.BlogEntries[j].ID
}
