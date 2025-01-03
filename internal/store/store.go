package store

import (
	"sync"
	"time"
)

type Entry struct {
	Value  string
	Expiry time.Time
}

func (e Entry) IsExpired() bool {
	if e.Expiry.IsZero() {
		return false
	}
	return time.Now().After(e.Expiry)
}

func NewEntry(val string, expirationDuration time.Duration) Entry {
	expirationDate := time.Time{}
	if expirationDuration != 0 {
		expirationDate = time.Now().Add(expirationDuration)
	}

	return Entry{
		Value:  val,
		Expiry: expirationDate,
	}
}

type DB = map[string]*Entry

type Store struct {
	db DB
	mu sync.Mutex
}

func NewStore() *Store {
	return &Store{
		db: make(DB),
	}
}

func (store *Store) Set(key, value string, expiry time.Duration) {
	store.mu.Lock()
	defer store.mu.Unlock()

	entry := NewEntry(value, expiry)
	store.db[key] = &entry
}

func (store *Store) Get(key string) string {
	entry := store.db[key]
	if entry == nil {
		return ""
	}

	if entry.IsExpired() {
		return ""
	}

	return entry.Value
}
