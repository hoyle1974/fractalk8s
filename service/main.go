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

func mandelIteration3(bcx, bcy *big.Float, maxIter int) (*big.Float, int) {
	// func Iter(p complex128, maxit uint) uint {
	// z := p
	p_real, _ := bcx.Float64()
	p_imag, _ := bcy.Float64()

	z_real := p_real
	z_imag := p_imag

	for it := 0; it < maxIter; it++ {
		//z = z*z + p
		temp_real := z_real*z_real - z_imag*z_imag
		temp_imag := z_real*z_imag + z_real*z_imag
		z_real = temp_real + p_real
		z_imag = temp_imag + p_imag

		// For better performance we use r^2 + i^2 > 4 instead of cmplx.Abs(z) > 2
		// if r, i := real(z), imag(z); r*r+i*i > 4 {
		if z_real*z_real+z_imag*z_imag > 4 {
			return big.NewFloat(0), it
		}
	}

	return big.NewFloat(0), maxIter

}

func mandelIteration1(bcx, bcy *big.Float, maxIter int) (*big.Float, int) {
	cx, _ := bcx.Float64()
	cy, _ := bcy.Float64()

	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return big.NewFloat(xx + yy), i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}

	logZn := (x*x + y*y) / 2
	return big.NewFloat(logZn), maxIter
}

func mandelIteration2(cx, cy *big.Float, maxIter int) (*big.Float, int) {

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
	cy, _, err = cy.Parse(r.URL.Query().Get("y"), 10)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing y")
	}
	iter, err := strconv.Atoi(r.URL.Query().Get("i"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing i")
	}

	value, iter := mandelIteration2(cx, cy, iter)

	log.Info().Msgf("value=%v iter=%d", value.String(), iter)
	fmt.Fprintf(w, "%v\n%d\n", value, iter)
}

func main() {
	log.Info().Msg("Setup Handler")
	http.HandleFunc("/iter", iterhandler)

	log.Info().Msg("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Unexpected error")
}
