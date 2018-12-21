package db

import (
	"encoding/binary"
	"time"

	"github.com/golang/protobuf/proto"
	bolt "go.etcd.io/bbolt"
)

// Open the database at the given path.
func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

// DB wraps a Bolt database.
type DB struct {
	db *bolt.DB
}

// Close the database.
func (db *DB) Close() error {
	return db.db.Close()
}

func bucket(t time.Time) []byte {
	return []byte(t.Local().Format("2006-01-02"))
}

// AddRecord adds a record to the database.
func (db *DB) AddRecord(record *Record) error {
	value, err := proto.Marshal(record)
	if err != nil {
		return err
	}
	t := record.GetTime()
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket(*t))
		if err != nil {
			return err
		}
		var key [8]byte
		seq, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint64(key[:], seq)
		return bucket.Put(key[:], value)
	})
}

// History returns the history for the given day.
func (db *DB) History(day time.Time) ([]*Record, error) {
	var history []*Record
	err := db.db.View(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket(bucket(day))
		if bucket == nil {
			return nil
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var record Record
			if err = proto.Unmarshal(value, &record); err != nil {
				return err
			}
			history = append(history, &record)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return history, nil
}
