package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

/*
func mandelIteration(cx, cy *big.Float, maxIter int) (*big.Float, int) {
	var x, y, xx, yy big.Float

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}

	logZn := (x*x + y*y) / 2
	return logZn, maxIter
}
*/

var four = big.NewFloat(float64(4))
var two = big.NewFloat(float64(2))
var half = big.NewFloat(float64(0.5))

func mandelIteration(cx, cy *big.Float, maxIter int) (*big.Float, int) {

	x := big.NewFloat(0.0)
	y := big.NewFloat(0.0)
	xx := big.NewFloat(0.0)
	yy := big.NewFloat(0.0)
	xy := big.NewFloat(0.0)
	temp := big.NewFloat(0.0)
	logZn := big.NewFloat(0.0)

	for i := 0; i < maxIter; i++ {
		xy.Mul(x, y)
		xx.Mul(x, x)
		yy.Mul(y, y)

		temp.Add(xx, yy)
		if temp.Cmp(four) > 0 {
			return temp, i
		}
		// x = xx - yy + cx
		x.Sub(xx, yy)
		x.Add(x, cx)
		// y = 2*xy + cy
		y.Mul(x, y)
		y.Mul(two, y)
		y.Add(y, cy)
	}

	//logZn := (x*x + y*y) / 2
	xx.Mul(x, x)
	yy.Mul(y, y)
	logZn.Add(xx, yy)
	logZn.Mul(logZn, half)

	return logZn, maxIter
}

func iterhandler(w http.ResponseWriter, r *http.Request) {

	cx := big.NewFloat(0.0)
	cy := big.NewFloat(0.0)

	log.Info().Msgf("x=%s y=%s i=%s\n", r.URL.Query().Get("x"), r.URL.Query().Get("y"), r.URL.Query().Get("i"))

	cx, _, err := cx.Parse(r.URL.Query().Get("x"), 10)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing x")
	}
	cy, _, err = cx.Parse(r.URL.Query().Get("y"), 10)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing y")
	}
	iter, err := strconv.Atoi(r.URL.Query().Get("i"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing i")
	}

	value, iter := mandelIteration(cx, cy, iter)

	fmt.Fprintf(w, "%v\n%d\n", value, iter)
}

func main() {
	log.Info().Msg("Setup Handler")
	http.HandleFunc("/iter", iterhandler)

	log.Info().Msg("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Unexpected error")
}
