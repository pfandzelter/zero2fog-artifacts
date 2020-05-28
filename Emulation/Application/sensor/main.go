package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var adaptEndpoint string = fmt.Sprintf("http://%s:%s/sensor", os.Getenv("ADAPT_IP"), os.Getenv("ADAPT_PORT"))

const interval int = 10
const upper int = 150
const lower int = 50

func main() {
	// update the world every 10 milliseconds
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			temp := rand.Intn(upper-lower) + lower
			// send random temperature

			id, err := uuid.NewRandom()
			if err != nil {
				log.Print(err)
				continue
			}

			type Reading struct {
				Temp      int    `json:"temp"`
				UUID      string `json:"uuid"`
				Timestamp string `json:"timestamp"`
			}

			data, err := json.Marshal(Reading{
				Temp:      temp,
				UUID:      id.String(),
				Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
			})

			if err == nil {
				log.Printf("send,adapt,%s,%s", id.String(),strconv.FormatInt(time.Now().UnixNano(), 10))
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
}
