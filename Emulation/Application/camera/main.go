package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var cfdEndpoint string = fmt.Sprintf("http://%s:%s/image", os.Getenv("CFD_IP"), os.Getenv("CFD_PORT"))

// update interval in milliseconds
const interval int = 100

const width int = 100
const height int = 100

func generateImage() string {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(2) == 1 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	var r bytes.Buffer

	err := png.Encode(&r, img)

	if err != nil {
		return ""
	}

	reader := bufio.NewReader(&r)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}

func main() {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	for {
		<-ticker.C
		img := generateImage()

		id, err := uuid.NewRandom()
		if err != nil {
			log.Print(err)
			continue
		}

		type Request struct {
			Img       string `json:"img"`
			UUID      string `json:"uuid"`
			Timestamp string `json:"timestamp"`
		}

		data, err := json.Marshal(Request{
			Img:       img,
			UUID:      id.String(),
			Timestamp: strconv.FormatInt(time.Now().UnixNano(), 10),
		})

		if err != nil {
			return
		}

		log.Printf("send,camera,%s,%s", id.String(), strconv.FormatInt(time.Now().UnixNano(), 10))

		req, err := http.NewRequest("POST", cfdEndpoint, bytes.NewReader(data))

		if err == nil {
			_, err := (&http.Client{}).Do(req)

			if err != nil {
				log.Print(err)
			}
		}

	}

}
