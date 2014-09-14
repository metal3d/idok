package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// when quiting (CTRL+C for example) - tell to XBMC to stop
func OnQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-c:
		fmt.Println("Quiting")
		resp := getActivePlayer()
		var playerid int
		for _, result := range resp.Result {
			for key, val := range result {
				if key == "playerid" {
					playerid = int(val.(float64))
				}
			}
		}

		http.Post(HOST, "application/json", bytes.NewBufferString(fmt.Sprintf(STOPBODY, playerid)))
		os.Exit(0)
	}
}
