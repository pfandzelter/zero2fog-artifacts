package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var prediction string

func main() {
	type Prediction struct {
		Prediction string `json:"prediction"`
		UUID       string `json:"uuid"`
		Timestamp  string `json:"timestamp"`
	}

	http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var data Prediction
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("recv,input,%s,%s,%s", data.UUID, data.Timestamp, timestamp)

		prediction = data.Prediction
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		res := Prediction{
			Prediction: prediction,
		}

		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			return
		}
	})

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PROGNOSIS_PORT"), nil))
}
