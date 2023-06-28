package main

import (
	"math/big"
	"net/http"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
	},
	Timeout: 0 * time.Second,
}

func NewWorkerRequest(x, y int, centerX, centerY, size *big.Float, out []byte) func() {
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

		it := calc(client, xx, yy, 128)

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
