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
	Timeout: 10 * time.Second,
}

func NewWorkerRequest(x, y int, centerX, centerY, size *big.Float, out []byte) func() {
	return func() {

		ti := 0
		xx := make([]*big.Float, chunk*chunk)
		yy := make([]*big.Float, chunk*chunk)

		for i := 0; i < chunk; i++ {
			for j := 0; j < chunk; j++ {
				// Center around scrren
				xx[ti] = big.NewFloat(float64(x + i - (screenWidth / 2)))
				yy[ti] = big.NewFloat(float64(y + j - (screenHeight / 2)))

				// scale
				xx[ti].Quo(xx[ti], big.NewFloat(screenWidth)) // -.5 to .5
				yy[ti].Quo(yy[ti], big.NewFloat(screenHeight))

				// scale
				xx[ti].Mul(xx[ti], size)
				yy[ti].Mul(yy[ti], size)

				// Offset
				xx[ti].Add(xx[ti], centerX)
				yy[ti].Add(yy[ti], centerY)

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
