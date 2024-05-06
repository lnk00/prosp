package db

import (
	"encoding/json"
	"log"

	"github.com/lnk00/prosp/models"
	bolt "go.etcd.io/bbolt"
)

type Database struct {
	Client *bolt.DB
}

func New() Database {
	db, err := bolt.Open("./prosp.db", 0600, nil)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("jobs"))
		return nil
	})

	return Database{
		Client: db,
	}
}

func (db Database) GetJobs() {
	db.Client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("jobs"))
		bucket.ForEach(func(k, v []byte) error {
			log.Printf("Key: %s\tvalue: %s", k, v)
			return nil
		})
		return nil
	})
}

func (db Database) SaveJob(job models.Job) {
	db.Client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("jobs"))
		encoded, err := json.Marshal(job)
		if err != nil {
			log.Printf("failed to encode job data: %v", err)
		}
		bucket.Put([]byte(job.Link), encoded)
		return nil
	})
}

func (db Database) SaveAllJobs(jobs []models.Job) {
	for _, job := range jobs {
		db.SaveJob(job)
	}
}
