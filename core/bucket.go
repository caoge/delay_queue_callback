package core

import (
	"delay_queue_callback/storage"
)

type BucketItem struct {
	JobSign string
	ExecTime int64
}

type Bucket struct {
	name string
}

func (bucket *Bucket) GetItem() (*BucketItem, error) {
	val, execTime, err := storage.GetTop(bucket.name)
	if err != nil {
		return nil, err
	}

	if val == "" {
		return nil, nil
	}


	return &BucketItem{val, execTime}, nil

}

func (bucket *Bucket) AddItem(jobSign string, delay int64) error {
	err := storage.Zadd(bucket.name, jobSign, delay)
	if err != nil {
		return err
	}

	return nil
}

func (bucket *Bucket) DelItem(jobSign string) error {
	err := storage.Zrem(bucket.name, jobSign)

	if err != nil {
		return err
	}

	return nil
}
