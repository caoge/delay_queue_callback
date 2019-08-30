package core

import (
	"delay_queue_callback/storage"
	"encoding/json"
	"errors"
	"fmt"
)

type Job struct {
	Topic    string `json:"topic"`
	Id       string `json:"id"`
	ExecTime int64  `json:"exec_time"`
	Body     interface{} `json:"body"`
	Callback string `json:"callback"`
	JobSign  string `json:"job_sign"`
}

func GetJob(jobSign string) (*Job, error) {
	val, err := storage.Get(jobSign)
	if err != nil {
		return nil, errors.New("从存储引擎中获取出错")
	}

	if val == "" {
		return nil, nil
	}

	job := new(Job)
	err = json.Unmarshal([]byte(val), job)
	if err != nil {
		return nil, errors.New("解析job json出错")
	}

	return job, err
}

func PushJob(job Job) error {
	jobJson, err := json.Marshal(job)
	if err != nil {
		return errors.New(fmt.Sprintf("json序列化出错%v", job))
	}

	err = storage.Set(job.JobSign, string(jobJson))
	if err != nil {
		return errors.New(fmt.Sprintf("job加入存储引擎出错%v", job))
	}

	return nil
}

func RemoveJob(jobSign string) error {
	err := storage.Del(jobSign)
	return err
}

func JobSign(topic, id string) string {
	return topic + "-" + id
}
