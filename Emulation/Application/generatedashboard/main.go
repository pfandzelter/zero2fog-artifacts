package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot/vg"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

var dashboardEndpoint string = fmt.Sprintf("http://%s:%s/input", os.Getenv("CENTRALDASHBOARD_IP"), os.Getenv("CENTRALDASHBOARD_PORT"))

// amount of historic data to save
const historic int = 1000

type PackCtrlData struct {
	Rate      int    `json:"rate"`
	Backlog   int    `json:"backlog"`
	UUID      string `json:"uuid"`
	Timestamp string `json:"timestamp"`
}

type Store struct {
	Data    []PackCtrlData
	Current int
	sync.Mutex
}

var s Store

func encode(data []int) (string, error) {
	p, err := plot.New()

	if err != nil {
		return "", err
	}

	p.Title.Text = "histogram plot"

	var valBacklog plotter.Values

	for _, d := range data {
		valBacklog = append(valBacklog, float64(d))
	}

	hist, err := plotter.NewHist(valBacklog, historic)

	if err != nil {
		return "", err
	}

	p.Add(hist)

	w, err := p.WriterTo(10*vg.Centimeter, 10*vg.Centimeter, "png")

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = w.WriteTo(buf)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func update(d PackCtrlData) {

	data := make([]PackCtrlData, historic)

	s.Lock()
	s.Current = (s.Current + 1) % historic

	s.Data[s.Current] = d

	c := s.Current
	copy(data, s.Data)
	s.Unlock()

	type Response struct {
		Backlog   string `json:"backlog"`
		Rate      string `json:"backlog"`
		UUID      string `json:"uuid"`
		Timestamp string `json:"timestamp"`
	}

	backlogData := make([]int, historic)
	rateData := make([]int, historic)

	for i, d := range data {
		offset := (c + i) % historic
		backlogData[offset] = d.Backlog
		rateData[offset] = d.Rate
	}

	encBacklog, err := encode(backlogData)

	if err != nil {
		log.Print(err)
		return
	}

	encRate, err := encode(rateData)

	if err != nil {
		log.Print(err)
		return
	}

	res := Response{
		Backlog:   encBacklog,
		Rate:      encRate,
		UUID:      d.UUID,
		Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	text, err := json.Marshal(res)

	if err != nil {
		log.Print(err)
		return
	}

	req, err := http.NewRequest("POST", dashboardEndpoint, bytes.NewReader(text))

	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("send,generate_dashboard,%s,%s", d.UUID, strconv.FormatInt(time.Now().UnixNano(), 10))

	_, err = (&http.Client{
		Timeout: 5 * time.Second,
	}).Do(req)

	if err != nil {
		log.Print(err)
		return
	}

}

func main() {
	// HTTP service, collects historic data as well and sends it out to a frontend if requested

	s = Store{
		Data:    make([]PackCtrlData, historic),
		Current: 0,
	}

	// fill Store with random data first
	t := time.Now().UnixNano()

	for i := range s.Data {
		s.Data[i] = PackCtrlData{
			Rate:      rand.Intn(1000),
			Backlog:   rand.Intn(1000),
			UUID:      "invalid",
			Timestamp: fmt.Sprintf("%d", rand.Int63n(t)),
		}
	}

	http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

		var d PackCtrlData
		err := json.NewDecoder(r.Body).Decode(&d)

		if err != nil {
			return
		}

		log.Printf("recv,input,%s,%s,%s", d.UUID, d.Timestamp, timestamp)

		go update(d)

	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("GENERATEDASHBOARD_PORT"), nil))
}
