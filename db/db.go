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

func (db Database) GetJobs() []models.Job {
	jobs := []models.Job{}

	db.Client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("jobs"))
		bucket.ForEach(func(k, v []byte) error {
			var job models.Job
			err := json.Unmarshal(v, &job)
			if err != nil {
				log.Fatalf("failed to decode job data: %v", err)
			}

			jobs = append(jobs, job)

			return nil
		})
		return nil
	})

	return jobs
}

func (db Database) SaveJob(job models.Job) {
	db.Client.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("jobs"))

		value := bucket.Get([]byte(job.Link))
		if value != nil {
			return nil
		}

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
