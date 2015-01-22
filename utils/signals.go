package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// when quiting (CTRL+C for example) - tell to XBMC to stop.
func OnQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	fmt.Println("Quiting")
	resp := getActivePlayer()
	var playerid int
	for _, result := range resp.Result {
		// TODO maybe result["playerid"] is ok...
		for key, val := range result {
			if key == "playerid" {
				playerid = int(val.(float64))
			}
		}
	}
	// tell to Kodi to stop
	http.Post(GlobalConfig.JsonRPC, "application/json", bytes.NewBufferString(fmt.Sprintf(STOPBODY, playerid)))
	os.Exit(0)
}
