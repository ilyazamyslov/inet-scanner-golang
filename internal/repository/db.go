package repository

import (
	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	riaken_core "github.com/riaken/riaken-core"
)

type DB struct {
	DB *riaken_core.Session
}

func NewDB(db *riaken_core.Session) *DB {
	return &DB{db}
}

func (db *DB) Load(key string) (model.Host, bool) {
	return model.Host{}, true
}

func (db *DB) Store(key string, value model.Host) {
	return
}
