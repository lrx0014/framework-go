package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
)

var useWebsockets = flag.Bool("websockets", false, "Whether to use websockets")

type Message struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Client.
func main() {
	flag.Parse()

	log.Println("Client Running...")

	if *useWebsockets {
		ws, err := websocket.Dial("ws://localhost:8080/", "", "http://localhost:8080")
		for {
			var m Message
			err = websocket.JSON.Receive(ws, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
				glog.Fatal(err)
			}
			glog.Infof("Received: %+v", m)
		}
	} else {
		glog.Info("Sending request...")
		req, err := http.NewRequest("GET", "http://localhost:8080", nil)
		if err != nil {
			glog.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			glog.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			glog.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
		}

		dec := json.NewDecoder(resp.Body)
		for {
			var m Message
			err := dec.Decode(&m)
			if err != nil {
				if err == io.EOF {
					break
				}
				glog.Fatal(err)
			}
			glog.Infof("Got response: %+v", m)
		}
	}

	glog.Infof("Server finished request...")
}
