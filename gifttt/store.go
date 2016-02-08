package gifttt

import (
	"errors"

	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/boltdb/bolt"
)

const (
	fileMode = 0600
)

var (
	_BUCKET        = []byte("gifttt")
	_store  *Store = nil

	ErrUnknownBucket = errors.New("bucket '" + string(_BUCKET) + "' does not exist")
	ErrNotFound      = errors.New("key not found")
)

type Store struct {
	db   *bolt.DB
	path string
}

func GetStore() *Store {
	return _store
}

func StoreInit(path string) error {
	// open the BoltDB and return error if this did not work
	// (beware that the same DB can only be opened by one
	// process)
	handle, err := bolt.Open(path, fileMode, nil)
	if err != nil {
		return err
	}

	// create the store
	store := &Store{
		db:   handle,
		path: path,
	}

	// initialize the database with needed buckets
	err = store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(_BUCKET)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		store.Close()
		return err
	}

	_store = store
	return nil
}

// close all connections to the underlying database file to be used
// by another process
func (store *Store) Close() {
	store.db.Close()
}

// set a key to the specified value.
func (store *Store) Set(key, value string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(_BUCKET)
		if b == nil {
			return ErrUnknownBucket
		}
		err := b.Put([]byte(key), []byte(value))

		return err
	})

	return err
}

// get the content of a key from the specified bucket
func (store *Store) Get(key string) (value string, err error) {
	err = store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(_BUCKET)
		if b == nil {
			return ErrUnknownBucket
		}
		data := b.Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		value = string(data)

		return nil
	})

	return value, err
}
