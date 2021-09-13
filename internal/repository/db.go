package repository

import (
	"encoding/json"
	"time"

	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	riaken_core "github.com/riaken/riaken-core"
	"github.com/rs/zerolog"
)

type DB struct {
	client *riaken_core.Client
	logger *zerolog.Logger
}

func NewDB(db *riaken_core.Client, logger *zerolog.Logger) *DB {
	return &DB{
		client: db,
		logger: logger,
	}
}

func (db *DB) Load(key string) (model.Host, bool) {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	object := bucket.Object(key)

	data, err := object.Fetch()

	if err != nil {
		return model.Host{}, false
	}
	if len(data.GetContent()) < 1 {
		return model.Host{}, false
	}
	jData := data.GetContent()[0].GetValue()

	var host model.Host
	err = json.Unmarshal(jData, &host)

	if err != nil {
		return model.Host{}, false
	}

	timeNow := time.Now().Unix()
	deadTime := time.Unix(host.Timestamp, 0).Add(24 * time.Hour).Unix()
	if timeNow > deadTime {
		db.logger.Info().Msg("Deleting " + key)
		if _, err := object.Delete(); err != nil {
			return model.Host{}, false
		}
		return model.Host{}, false
	}
	return host, true
}

func (db *DB) Store(key string, value model.Host) error {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	object := bucket.Object(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := object.Store(data); err != nil {
		return err
	}
	return nil
}
