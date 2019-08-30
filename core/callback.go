package core

import (
	"delay_queue_callback/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func call(url string, body io.Reader) {

	timeout := time.Duration(config.Conf.Timeout) * time.Second

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return
	}

	var content []byte

	for i := 0; i < config.Conf.MaxTries; i++ {
		res, err := client.Do(req)
		if err != nil {
			log.Printf("第%d次请求失败,失败原因为%v", i+1, err)
			continue
		}

		content, err = ioutil.ReadAll(res.Body)
		if string(content) == "ok" {
			return
		}else{
			log.Printf("第%d次请求失败,返回为%v", i+1, content)
		}

	}
}
