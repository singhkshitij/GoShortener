package store

import (
	"bytes"
	"log"
	"net/http"
	"strconv"

	"github.com/etcd-io/bbolt"
)

// Store is the store interface for urls.
type Store interface {
	Set(key string, value string) error
	Get(key string) string
	GetSize() int
	Close()
}

func panic(v interface{}) {
	log.Panic(v)
}

var tableURLs = []byte("urls")

// DB representation of a Store.
type DB struct {
	db *bbolt.DB
}

var _ Store = &DB{}

func openDBConnec(dBPath string) *bbolt.DB {

	db, err := bbolt.Open(dBPath, 0600, nil)
	if err != nil {
		panic(err)
	}

	tables := [...][]byte{
		tableURLs,
	}

	db.Update(func(tx *bbolt.Tx) (err error) {
		for _, bucket := range tables {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				panic(err)
			}
		}
		return
	})

	return db
}

// NewDB returns a new DB instance, its connection is opened.
func NewDB(dbpath string) *DB {
	return &DB{
		db: openDBConnec(dbpath),
	}
}

// Clear clears all the database entries for the table urls.
func (dbase *DB) Clear() error {
	return dbase.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(tableURLs)
	})
}

// Close shutdowns the data(base) connection.
func (dbase *DB) Close() {
	if err := dbase.db.Close(); err != nil {
		panic(err)
	}
}

// Set sets a shorten url and its key
func (dbase *DB) Set(key string, value string) error {
	return dbase.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tableURLs)
		if err != nil {
			panic(err)
		}

		cursor := bucket.Cursor()
		byteValue := []byte(value)

		for bucketKey, bucketValue := cursor.First(); bucketKey != nil; bucketKey, bucketValue = cursor.Next() {
			if bytes.Equal(bucketValue, byteValue) {
				bucket.Delete(bucketKey)
				break
			}
		}
		return bucket.Put([]byte(key), byteValue)
	})
}

// Get long url from short id or an empty string if not found.
func (dbase *DB) Get(key string) (value string) {
	byteKey := []byte(key)
	dbase.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(tableURLs)
		if bucket == nil {
			return nil
		}

		byteValue := bucket.Get(byteKey)
		if byteValue != nil {
			value = string(byteValue)
		}
		return nil
	})
	return
}

// GetSize all the "shorted" urls length
func (dbase *DB) GetSize() (size int) {
	dbase.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(tableURLs)
		if bucket == nil {
			return nil
		}

		size = bucket.Stats().KeyN
		return nil
	})
	return
}

//Backup database to provided filename
func (dbase *DB) Backup(bkpFileName string, writer http.ResponseWriter) http.ResponseWriter {
	dbase.db.View(func(tx *bbolt.Tx) error {
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Disposition", `attachment; filename="`+bkpFileName+`"`)
		writer.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(writer)
		if err!= nil {
			panic(err)
		}
		return nil
	})
	return writer
}
