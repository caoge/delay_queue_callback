package main

import (
	"delay_queue_callback/config"
	"delay_queue_callback/routes"
	"fmt"
	"log"
	"net/http"
)

type Cmd struct{}

func main() {
	//fp, _ := os.OpenFile(config.Conf.LogPath, os.O_WRONLY|os.O_APPEND|os.O_SYNC|os.O_CREATE,
	//	0755)
	//log.SetOutput(fp)
	http.HandleFunc("/job",routes.Handle)

	fmt.Printf("listening %s\n", config.Conf.Listen)

	err := http.ListenAndServe(config.Conf.Listen,nil)
	if err != nil {
		log.Fatal(err)
	}

}