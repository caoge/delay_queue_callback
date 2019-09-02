package core

import (
	"delay_queue_callback/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var bucket <-chan *Bucket

func init() {
	bucket = bucketGenerator()
	initTimers()
}

func initTimers() {
	var ticker *time.Ticker

	for i := 0; i < config.Conf.BucketSize; i++ {
		ticker = time.NewTicker(1 * time.Second)
		go scanBucket(ticker, <-bucket)
	}
}

func scanBucket(ticker *time.Ticker, bk *Bucket) {
	for {
		select {
		case <-ticker.C:
			scanHandle(bk)
		}
	}

}

func bucketGenerator() <-chan *Bucket {
	var bk *Bucket

	var num = 1

	var bucketChannel = make(chan *Bucket)

	go func() {
		for {
			bk = new(Bucket)
			bk.name = fmt.Sprintf(config.Conf.BucketName, num)

			bucketChannel <- bk

			if num == config.Conf.BucketSize {
				num = 1
			} else {
				num++
			}
		}
	}()

	return bucketChannel
}

func scanHandle(bk *Bucket) {
	for {
		bucketItem, err := bk.GetItem()
		if err != nil {
			log.Println("获取bucket出错")
			return
		}

		if bucketItem == nil {
			return
		}

		if bucketItem.ExecTime <= time.Now().Unix() {
			job, err := GetJob(bucketItem.JobSign)
			if err != nil {
				clear(bk, bucketItem.JobSign)
				log.Println(err)
				continue
			}

			//说明job被删除
			if job == nil {
				bk.DelItem(bucketItem.JobSign)
				continue
			}

			jsonBody, err := json.Marshal(job) //todo getJob unmarshal了这里marshal好像有点傻
			if err != nil {
				log.Println("json序列化失败", job)
			}
			clear(bk, bucketItem.JobSign)
			log.Println(string(jsonBody))
			go call(job.Callback, jsonBody)
		} else {
			return
		}
	}

}

func NewWork(job Job) error {
	err := PushJob(job)
	if err != nil {
		return err
	}

	err = (<-bucket).AddItem(job.JobSign, job.Delay+time.Now().Unix())
	return err

}

func clear(bk *Bucket, jobSign string) {
	RemoveJob(jobSign)
	bk.DelItem(jobSign)
}
