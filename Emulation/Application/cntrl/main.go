package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var adaptEndpoint string = fmt.Sprintf("http://%s:%s/prodcntrl", os.Getenv("ADAPT_IP"), os.Getenv("ADAPT_PORT"))

// standard production rate
const rate = 100

// update interval in milliseconds
const interval int = 100

func update(queue <-chan struct{}) {
	// update the world every 100 milliseconds
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	discarded := 0
	for {
		select {
		case <-queue:
			discarded++
		case <-ticker.C:
			// send rate - discarded

			curr := (rate * (interval / 1000)) - discarded

			id, err := uuid.NewRandom()
			if err != nil {
				log.Print(err)
				continue
			}

			type ProdCntrlData struct {
				ProdRate  int    `json:"prod_rate"`
				UUID      string `json:"uuid"`
				Timestamp string `json:"timestamp"`
			}

			data, err := json.Marshal(ProdCntrlData{
				ProdRate:  curr,
				UUID:      id.String(),
				Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			})

			if err != nil {
				return
			}
			log.Printf("send,prodctrl,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

			req, err := http.NewRequest("POST", adaptEndpoint, bytes.NewReader(data))

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
	discardqueue := make(chan struct{})

	type Request struct {
		UUID      string `json:"uuid"`
		Timestamp string `json:"timestamp"`
	}

	http.HandleFunc("/discard", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data Request
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			return
		}

		log.Printf("recv,discard,%s,%s,%s", data.UUID, data.Timestamp, timestamp)

		discardqueue <- struct{}{}
	})

	go update(discardqueue)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("CNTRL_PORT"), nil))

}
