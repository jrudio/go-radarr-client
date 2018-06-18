package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

var isVerbose bool

type store struct {
	db       *badger.DB
	isClosed bool
	keys     storeKeys
	secret   []byte
}

type storeKeys struct {
	appSecret []byte
	radarrKey []byte
	radarrURL []byte
}

func initDataStore(dirName string) (store, error) {
	var db store

	if isVerbose {
		fmt.Println("checking if our database exists in the home directory at:", dirName)
	}

	// create a directory for our database
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if isVerbose {
			fmt.Println("creating directory because it doesn't exist")
		}

		if err := os.Mkdir(dirName, os.ModePerm); err != nil {
			return db, err
		}
	}

	options := badger.DefaultOptions

	options.Dir = dirName
	options.ValueDir = dirName

	kvStore, err := badger.Open(options)

	if err != nil {
		return db, err
	}

	if isVerbose {
		fmt.Println("successfully opened data store")
	}

	db.db = kvStore
	db.keys = storeKeys{
		radarrKey: []byte("radarr-key"),
		radarrURL: []byte("radarr-url"),
	}

	return db, nil
}

func (s store) Close() {
	if s.isClosed {
		fmt.Println("data store already closed")
		return
	}

	if err := s.db.Close(); err != nil {
		fmt.Printf("data store failed to closed: %v\n", err)
	}

	s.isClosed = true
}

func (s store) getRadarrKey() (string, error) {
	var radarrKey string

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.radarrKey)

		if err != nil {
			return err
		}

		_radarrKey, err := item.Value()

		if err != nil {
			return err
		}

		radarrKey = string(_radarrKey)

		return nil
	}); err != nil {
		return radarrKey, err
	}

	if isVerbose {
		fmt.Printf("Your radarr key is %s\n", radarrKey)
	}

	return radarrKey, nil
}

func (s store) saveRadarrKey(key string) error {
	if isVerbose {
		fmt.Printf("your radarr key: %s\n", string(key))
	}

	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.radarrKey, []byte(key), 0x00)
	}); err != nil {
		return err
	}

	if isVerbose {
		fmt.Println("saved key to store")
	}

	return nil
}

func (s store) getRadarrURL() (string, error) {
	var radarrURL string

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.radarrURL)

		if err != nil {
			return err
		}

		serializedServer, err := item.Value()

		if err != nil {
			return err
		}

		radarrURL = string(serializedServer)

		return nil
	})

	return radarrURL, err
}

func (s store) saveRadarrURL(radarrURL string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.radarrURL, []byte(radarrURL), 0)
	})
}
