package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var store struct {
	rate    string
	backlog string
	sync.Mutex
}

func update(rate, backlog string) {
	store.Lock()

	store.rate = rate
	store.backlog = backlog

	store.Unlock()
}

func main() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {

		id, err := uuid.NewRandom()

		if err != nil {
			log.Print(err)
			return
		}

		type Response struct {
			Backlog   string `json:"backlog"`
			Rate      string `json:"backlog"`
			UUID      string `json:"uuid"`
			Timestamp string `json:"timestamp"`
		}

		res := Response{
			Rate:      store.rate,
			Backlog:   store.backlog,
			UUID:      id.String(),
			Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		}

		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			return
		}
	})

	http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		type Request struct {
			Backlog   string `json:"backlog"`
			Rate      string `json:"backlog"`
			UUID      string `json:"uuid"`
			Timestamp string `json:"timestamp"`
		}

		var d Request

		err := json.NewDecoder(r.Body).Decode(&d)

		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("recv,input,%s,%s,%s", d.UUID, d.Timestamp, timestamp)

		go update(d.Rate, d.Backlog)
	})

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("CENTRALDASHBOARD_PORT"), nil))
}
