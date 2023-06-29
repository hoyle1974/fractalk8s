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

func GetClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
		},
		Timeout: 60 * time.Second,
	}
}

func calc(client *http.Client, x, y []*big.Float, iter int) ([]int, time.Duration, time.Duration, int, int) {

	start := time.Now()

	req := common.NewMRequest(x, y, iter)

	// url := fmt.Sprintf("http://fractalk8s.decepticons.local/iter")
	url := "http://localhost:8080/iter"

	reqString := req.ToJsonString()
	reqSize := len(reqString)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqString)))
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

	response := common.NewMResponseFromJson(string(body))
	respSize := len(body)

	duration := time.Since(start)

	return response.Iter, response.CalcTime, duration, reqSize, respSize
}

func NewWorkerRequest(x, y int, centerX, centerY, size *big.Float, out []byte, metrics *Metrics) func() {
	return func() {

		ti := 0
		xx := make([]*big.Float, chunk*chunk)
		yy := make([]*big.Float, chunk*chunk)

		for i := 0; i < chunk; i++ {
			for j := 0; j < chunk; j++ {
				// Center around scrren
				// txx := big.NewRat(int64(x+i-(screenWidth/2)), screenWidth)
				txx := big.NewFloat(float64(x+i-(screenWidth/2)) / float64(screenWidth))

				// tyy := big.NewRat(int64(y+j-(screenHeight/2)), screenHeight)
				tyy := big.NewFloat(float64(y+j-(screenHeight/2)) / float64(screenHeight))

				// scale
				txx = big.NewFloat(0).SetPrec(0).Mul(txx, size)
				tyy = big.NewFloat(0).SetPrec(0).Mul(tyy, size)

				// Offset
				txx.Add(txx, centerX)
				tyy.Add(tyy, centerY)

				xx[ti] = txx
				yy[ti] = tyy

				ti++
			}
		}

		it, cd, rd, reqs, resps := calc(client, xx, yy, 128)
		metrics.AddDuration("calc", cd)
		metrics.AddDuration("request", rd)
		metrics.AddBytes("requests", int64(reqs))
		metrics.AddBytes("response", int64(resps))

		ti = 0
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
