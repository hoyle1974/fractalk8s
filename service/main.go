package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/hoyle1974/fractalk8s/common"
	"github.com/hoyle1974/fractalk8s/mandelbrot"
	"github.com/rs/zerolog/log"
)

func iterhandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Version 1")

	start := time.Now()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading body")
	}

	request := common.NewMRequestFromBytes(body)
	//cx, cy, iter := request.Extract()

	//for idx, _ := range cx {
	//	resultIter[idx] = mandelbrot.MandelbrotFloat(cx[idx], cy[idx], iter)
	//}

	chunk := request.Chunk
	x := request.X
	y := request.Y
	screenWidth := request.ScreenWidth
	screenHeight := request.ScreenHeight
	iter := request.Iter
	// centerX, centerY, size := request.ExtractFloats()
	centerX := request.CenterX
	centerY := request.CenterY
	size := request.Size

	centerX.SetPrec(centerX.MinPrec())
	centerY.SetPrec(centerY.MinPrec())
	size.SetPrec(size.MinPrec())

	resultIter := make([]int, chunk*chunk)

	small := size.Cmp(new(big.Float).SetFloat64(0.00000000000008)) > 0

	idx := 0
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

			if small {
				fx, _ := txx.Float64()
				fy, _ := tyy.Float64()
				resultIter[idx] = mandelbrot.MandelbrotFast(fx, fy, iter)
			} else {
				resultIter[idx] = mandelbrot.MandelbrotFloat(txx, tyy, iter)
			}
			idx++

		}
	}

	end := time.Now()
	calcTime := end.Sub(start)
	response := common.NewMResponse(resultIter, calcTime)

	fmt.Fprintf(w, "%s", response.ToBytes())
}

func main() {

	log.Info().Msg("Setup Handler")
	http.HandleFunc("/iter", iterhandler)

	log.Info().Msg("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Unexpected error")
}
