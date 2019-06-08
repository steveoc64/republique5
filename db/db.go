package db

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/memdebug"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DB is the controller for all BoltDB operations
type DB struct {
	log      *logrus.Logger
	db       *bolt.DB
	filename string
}

// Error declarations
var (
	ErrNoGameBucket = errors.New("No game bucket")
	ErrPutData      = errors.New("Put data")
	ErrProtoMarshal = errors.New("Marshal data")
)

// OpenReadDB opens a BoltDB for reading
func OpenReadDB(log *logrus.Logger, name string) (*DB, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)
	if _, err := os.Stat(filename); err == os.ErrExist {
		return nil, err
	}

	// open the DB
	db, err := bolt.Open(filename, 0644, &bolt.Options{ReadOnly: true, Timeout: 200 * time.Millisecond})
	if err != nil {
		if err == bolt.ErrTimeout {
			log.Fatal("DB is locked by another process - maybe republique serve is already running ?")
		}
		return nil, err
	}
	return &DB{
		log:      log,
		db:       db,
		filename: filename,
	}, nil
}

// OpenDB opens a BoltDB for writing
func OpenDB(log *logrus.Logger, name string) (*DB, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)
	if _, err := os.Stat(filename); err == os.ErrExist {
		return nil, err
	}

	// open the DB
	db, err := bolt.Open(filename, 0644, &bolt.Options{Timeout: 200 * time.Millisecond})
	if err != nil {
		if err == bolt.ErrTimeout {
			log.Fatal("DB already in use by another process")
		}
		return nil, err
	}
	return &DB{
		log:      log,
		db:       db,
		filename: filename,
	}, nil
}

// NewDB returns a new instance of a DB
func NewDB(log *logrus.Logger, name string) *DB {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)

	// Create the dir to store it all
	err := os.Mkdir(filepath.Dir(filename), 0744)
	if err != nil && !strings.HasSuffix(err.Error(), "file exists") {
		log.WithError(err).WithField("filename", filename).Warn("mkdir")
	}

	// Create the DB
	db, err := bolt.Open(filename, 0644, &bolt.Options{Timeout: 200 * time.Millisecond})
	if err != nil {
		if err == bolt.ErrTimeout {
			log.Fatal("DB already in use by another process")
		}
		log.Fatal(err)
	}

	// Create our game bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("game"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &DB{
		log:      log,
		db:       db,
		filename: filename,
	}
}

// Save writes a record to the DB
func (store *DB) Save(bucket, key string, data proto.Message) error {
	t1 := time.Now()
	err := store.db.Update(func(tx *bolt.Tx) error {
		game := tx.Bucket([]byte(bucket))
		if game == nil {
			store.log.Error(ErrNoGameBucket)
			return ErrNoGameBucket
		}
		b, err := proto.Marshal(data)
		if err != nil {
			store.log.WithError(err).Error(ErrProtoMarshal)
			return err
		}
		memdebug.Print(t1, "marshalled", bucket, key)
		err = game.Put([]byte(key), b)
		if err != nil {
			store.log.WithError(err).Error(ErrPutData)
			return err
		}
		return nil
	})
	if err != nil {
		store.log.WithError(err).WithField("filename", store.filename).Error(ErrPutData)
	}
	memdebug.Print(t1, "saved", bucket, key)
	return err
}

// Load loads a record from the DB
func (store *DB) Load(bucket, key string, data proto.Message) error {
	// retrieve the data
	t1 := time.Now()
	return store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %v not found!", b)
		}

		val := b.Get([]byte(key))
		memdebug.Print(t1, "loaded", bucket, key)
		err := proto.Unmarshal(val, data)
		memdebug.Print(t1, "unmarshal")
		return err
	})
}

// Close closes off a DB
func (store *DB) Close() error {
	return store.db.Close()
}
