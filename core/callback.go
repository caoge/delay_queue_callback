package core

import (
	"bytes"
	"delay_queue_callback/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func call(url string, body []byte, callSign chan<- bool) {

	timeout := time.Duration(config.Conf.Timeout) * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	var content []byte

	for i := 0; i < config.Conf.MaxTries; i++ {
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))

		defer req.Body.Close()
		if err != nil {
			log.Println("构建请求失败%v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			log.Printf("第%d次请求失败,失败原因为%v", i+1, err)
			continue
		}
		defer res.Body.Close()

		content, err = ioutil.ReadAll(res.Body)
		if string(content) == "ok" {
			callSign <- true
			return
		} else {
			callSign <- false
			log.Printf("第%d次请求失败,返回为%v", i+1, content)
		}

	}
}
