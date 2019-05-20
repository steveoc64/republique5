package republique

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type DB struct {
	log        *logrus.Logger
	db         *bolt.DB
	gameBucket []byte
	filename   string
}

var (
	ErrNoGameBucket = errors.New("No game bucket")
	ErrPutData      = errors.New("Put data")
	ErrProtoMarshal = errors.New("Marshal data")
)

func NewDB(log *logrus.Logger, name string) *DB {
	filename := filepath.Join(os.Getenv("HOME"), ".republique", name)

	// Create the dir to store it all
	err := os.Mkdir(filepath.Dir(filename), 0744)
	if err != nil {
		panic(err)
	}

	// Create the DB
	db, err := bolt.Open(filename, 0644, nil)
	if err != nil {
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
		log:        log,
		db:         db,
		gameBucket: []byte("game"),
		filename:   filename,
	}
}

func (store *DB) Save(data proto.Message) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		game := tx.Bucket(store.gameBucket)
		if game == nil {
			store.log.Error(ErrNoGameBucket)
			return ErrNoGameBucket
		}
		b, err := proto.Marshal(data)
		if err != nil {
			store.log.WithError(err).Error(ErrProtoMarshal)
			return err
		}
		err = game.Put([]byte("data"), b)
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

func (store *DB) Close() error {
	return store.db.Close()
}
