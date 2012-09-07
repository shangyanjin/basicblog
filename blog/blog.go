// Package blog contains blog and blog entry types and functions for loading,
// saving and using blogs
package blog

import (
	"encoding/json"
	"os"
	"time"
)

// Blog represents a blog with all of the entries belonging to it.
type Blog struct {
	Entries []BlogEntry
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
	if len(b.Entries) > 0 {
		e.ID = b.Entries[0].ID + 1
		b.Entries = append([]BlogEntry{*e}, b.Entries...)
	} else {
		e.ID = 1
		b.Entries = append(b.Entries, *e)
	}
}
