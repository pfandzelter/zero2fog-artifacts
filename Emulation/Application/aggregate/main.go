package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

const interval int = 10000

var generateDashboardEndpoint string = fmt.Sprintf("http://%s:%s/input", os.Getenv("GENERATEDASHBOARD_IP"), os.Getenv("GENERATEDASHBOARD_PORT"))

type PackCtrlData struct {
	Rate      int    `json:"rate"`
	Backlog   int    `json:"backlog"`
	UUID      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
}

func update(incoming <-chan PackCtrlData) {

	var avgbacklog float64

	var avgpackagingspeed float64

	var amount int

	var id string

	var timestamp string

	// update the world every interval seconds
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case data := <-incoming:
			avgbacklog = ((avgbacklog * float64(amount)) + float64(data.Backlog)) / (float64(amount) + 1)
			avgpackagingspeed = ((avgpackagingspeed * float64(amount)) + float64(data.Rate)) / (float64(amount) + 1)
			amount++
			id = data.UUID
			timestamp = strconv.FormatInt(time.Now().UnixNano(), 10)

		case <-ticker.C:
			// send data
			data, err := json.Marshal(PackCtrlData{
				Rate:      int(math.Round(avgpackagingspeed)),
				Backlog:   int(math.Round(avgbacklog)),
				UUID:      id,
				Timestamp: timestamp,
			})

			if err != nil {
				return
			}

			log.Printf("send,aggregate,%s,%s", id, timestamp)
			req, err := http.NewRequest("POST", generateDashboardEndpoint, bytes.NewReader(data))

			if err == nil {
				_, err := (&http.Client{}).Do(req)

				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func main() {
	// HTTP service, collects data and sends it out in 10 second intervals

	incoming := make(chan PackCtrlData)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data PackCtrlData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,data,%s,%s,%s", data.UUID, data.Timestamp, timestamp)

		incoming <- data
	})

	go update(incoming)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("AGGREGATE_PORT"), nil))
}
