package storage

import (
	"fmt"
)

var ErrAlreadyExists = fmt.Errorf("unable to create, entry already exists")
var ErrNotFound = fmt.Errorf("entry was not found")

// Store is a disk backed build store, holds the current status of builds
// and their outputs
type Store struct {
	dir    string
	memory map[string]Collection
}

// New creates and returns a Store that allows entries to be created,
//  appended, and retreived.
func New(dir string) (*Store, error) {
	// need to load it from disk.
	return &Store{
		dir:    dir,
		memory: map[string]Collection{},
	}, nil
}

// Create a new entry in the store, returning an error if one
// already exists
func (s *Store) Create(key string, val interface{}) error {
	_, found := s.memory[key]
	if found {
		return ErrAlreadyExists
	}
	s.memory[key] = []interface{}{val}
	return nil
}

// Append an entry to the collection, returning an error if it doesn't
// exist
func (s *Store) Append(key string, val interface{}) error {
	c, found := s.memory[key]
	if !found {
		return ErrNotFound
	}
	s.memory[key] = append(c, val)
	return nil
}

// Get a collection of entries from the store, return an error if none
// exist
func (s *Store) Get(key string) (Collection, error) {
	c, found := s.memory[key]
	if !found {
		return nil, ErrNotFound
	}
	return c, nil
}
