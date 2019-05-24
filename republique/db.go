package republique

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DB struct {
	log      *logrus.Logger
	db       *bolt.DB
	filename string
}

var (
	ErrNoGameBucket = errors.New("No game bucket")
	ErrPutData      = errors.New("Put data")
	ErrProtoMarshal = errors.New("Marshal data")
	gameBucket      = []byte("game")
	gameState       = []byte("state")
)

func OpenReadDB(log *logrus.Logger, name string) (*DB, error) {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)
	if _, err := os.Stat(filename); err == os.ErrExist {
		return nil, err
	}

	// open the DB
	db, err := bolt.Open(filename, 0644, &bolt.Options{ReadOnly: true})
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

func NewDB(log *logrus.Logger, name string) *DB {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)

	// Create the dir to store it all
	err := os.Mkdir(filepath.Dir(filename), 0744)
	if err != nil && !strings.HasSuffix(err.Error(), "file exists") {
		log.WithError(err).WithField("filename", filename).Warn("mkdir")
	}

	// Create the DB
	db, err := bolt.Open(filename, 0644,&bolt.Options{Timeout: 200 * time.Millisecond})
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

func (store *DB) Save(data proto.Message) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		game := tx.Bucket(gameBucket)
		if game == nil {
			store.log.Error(ErrNoGameBucket)
			return ErrNoGameBucket
		}
		b, err := proto.Marshal(data)
		if err != nil {
			store.log.WithError(err).Error(ErrProtoMarshal)
			return err
		}
		err = game.Put(gameState, b)
		if err != nil {
			store.log.WithError(err).Error(ErrPutData)
			return err
		}
		return nil
	})
	if err != nil {
		store.log.WithError(err).WithField("filename", store.filename).Error(ErrPutData)
	}
	return err
}

func (store *DB) Load(data proto.Message) error {
	// retrieve the data
	return store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(gameBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %v not found!", string(gameBucket))
		}

		val := bucket.Get(gameState)
		return proto.Unmarshal(val, data)
	})
}

func (store *DB) Close() error {
	return store.db.Close()
}
