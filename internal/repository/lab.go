package repository

import (
	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
)

func (db *DB) Query1(start string, end string) []model.Host {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	var keys [][]byte
	for out, err := bucket.ListKeys(); !out.GetDone(); out, err = bucket.ListKeys() {
		if err != nil {
			db.logger.Error().Err(err).Msg("ScanHost method error")
			break
		}
		keys = append(keys, out.GetKeys()...)
	}
	var keysToGet []string
	for _, val := range keys {
		if string(val) >= start || string(val) < end {
			keysToGet = append(keysToGet, string(val))
		}
	}
	var result []model.Host
	for _, val := range keysToGet {
		if host, ok := db.Load(val); ok {
			result = append(result, host)
		}
	}
	return result
}

//name not null
func (db *DB) Query2(start string, end string) []model.Host {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	var keys [][]byte
	for out, err := bucket.ListKeys(); !out.GetDone(); out, err = bucket.ListKeys() {
		if err != nil {
			db.logger.Error().Err(err).Msg("ScanHost method error")
			break
		}
		keys = append(keys, out.GetKeys()...)
	}
	var keysToGet []string
	for _, val := range keys {
		if string(val) >= start || string(val) < end {
			keysToGet = append(keysToGet, string(val))
		}
	}
	var result []model.Host
	for _, val := range keysToGet {
		if host, ok := db.Load(val); ok {
			if host.Name != "" {
				result = append(result, host)
			}
		}
	}
	return result
}

//os not null
func (db *DB) Query3(start string, end string) []model.Host {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	var keys [][]byte
	for out, err := bucket.ListKeys(); !out.GetDone(); out, err = bucket.ListKeys() {
		if err != nil {
			db.logger.Error().Err(err).Msg("ScanHost method error")
			break
		}
		keys = append(keys, out.GetKeys()...)
	}
	var keysToGet []string
	for _, val := range keys {
		if string(val) >= start || string(val) < end {
			keysToGet = append(keysToGet, string(val))
		}
	}
	var result []model.Host
	for _, val := range keysToGet {
		if host, ok := db.Load(val); ok {
			if host.Os != "" && len(host.Ports) > 0 {
				result = append(result, host)
			}

		}
	}
	return result
}

func (db *DB) Query4(start string, end string) int {
	session := db.client.Session()
	defer session.Release()
	bucket := session.GetBucket("hosts")
	var keys [][]byte
	for out, err := bucket.ListKeys(); !out.GetDone(); out, err = bucket.ListKeys() {
		if err != nil {
			db.logger.Error().Err(err).Msg("ScanHost method error")
			break
		}
		keys = append(keys, out.GetKeys()...)
	}
	var result int
	for _, val := range keys {
		if string(val) >= start || string(val) < end {
			result++
		}
	}
	return result
}
