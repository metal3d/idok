package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const TICK_CHECK = 1

// response of get players
type itemresp struct {
	Id      int
	Jsonrpc string
	Result  []map[string]interface{}
}

// begin to locally listen http to serve media
func Send(scheme, host, file string, port int) {

	u := url.URL{Path: file}
	file = u.String()
	addr := fmt.Sprintf("%s://%s:%d/%s", scheme, host, port, file)

	r, err := http.Post(HOST, "application/json", bytes.NewBufferString(fmt.Sprintf(BODY, addr)))
	if err != nil {
		log.Fatal(err)
	}
	response, _ := ioutil.ReadAll(r.Body)
	log.Println(string(response))
	// and wait media end
	go checkPlaying()
}

// send basic stream...
func SendBasicStream(uri string, local bool) {
	_body := fmt.Sprintf(BODY, uri)

	r, err := http.Post(HOST, "application/json", bytes.NewBufferString(_body))
	if err != nil {
		log.Fatal(err)
	}
	response, _ := ioutil.ReadAll(r.Body)
	log.Println(string(response))

	// handle CTRL+C to stop
	go OnQuit()

	// stay alive
	c := make(chan int)
	<-c
}

// Ask to play youtube video
func PlayYoutube(vidid string) {

	r, err := http.Post(HOST, "application/json", bytes.NewBufferString(fmt.Sprintf(YOUTUBEAPI, vidid)))
	if err != nil {
		log.Fatal(err)
	}
	response, _ := ioutil.ReadAll(r.Body)
	log.Println(string(response))

	// handle CTRL+C to stop
	go OnQuit()

	// stay alive
	c := make(chan int)
	<-c
}

// test if media is playing, if not then quit
func checkPlaying() {
	tick := time.Tick(TICK_CHECK * time.Second)
	for _ = range tick {
		resp := getActivePlayer()
		if len(resp.Result) == 0 {
			os.Exit(0)
		}
	}

}

// return active player from XBMC
func getActivePlayer() *itemresp {
	r, _ := http.Post(HOST, "application/json", bytes.NewBufferString(GETPLAYERBODY))
	response, _ := ioutil.ReadAll(r.Body)
	resp := &itemresp{}
	resp.Result = make([]map[string]interface{}, 0)
	json.Unmarshal(response, resp)
	return resp
}
