package main

import (
	"bytes"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/hoyle1974/fractalk8s/common"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        1,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
	},
	Timeout: 60 * time.Second,
}

func POST(client *http.Client, req common.MRequest) ([]byte, time.Duration, time.Duration, int, int) {
	start := time.Now()

	url := "http://localhost:8080/iter"
	reqBytes := req.ToBytes()
	reqSize := len(reqBytes)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	response := common.NewMResponseFromBytes(body)
	respSize := len(body)

	duration := time.Since(start)

	return response.Iter, response.CalcTime, duration, reqSize, respSize
}

func NewWorkerRequest(x, y int, centerX, centerY, size *big.Float, out []byte, metrics *Metrics) func() {
	return func() {

		request := common.NewMRequest(x, y, chunk, screenWidth, screenHeight, centerX, centerY, size, 255)
		it, cd, rd, reqs, resps := POST(client, request)

		metrics.AddDuration("calc", cd)
		metrics.AddDuration("request", rd)
		metrics.AddBytes("requests", int64(reqs))
		metrics.AddBytes("response", int64(resps))

		ti := 0
		for i := 0; i < chunk; i++ {
			for j := 0; j < chunk; j++ {
				r, g, b := color(it[ti])
				p := 4 * (((y + j) * screenWidth) + (x + i))
				out[p] = r
				out[p+1] = g
				out[p+2] = b
				out[p+3] = 0xff
				ti++
			}
		}

	}
}
