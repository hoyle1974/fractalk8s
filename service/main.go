package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/fracktalk8s/mandelbrot"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var nodeId = uuid.New()

/*
var four = big.NewFloat(float64(4))
var two = big.NewFloat(float64(2))
var half = big.NewFloat(float64(0.5))

func mandelIteration2(cx, cy *big.Float, maxIter int) (*big.Float, int) {

	x := big.NewFloat(0.0)
	y := big.NewFloat(0.0)
	xx := big.NewFloat(0.0)
	yy := big.NewFloat(0.0)
	// xy := big.NewFloat(0.0)
	temp := big.NewFloat(0.0)
	logZn := big.NewFloat(0.0)

	for i := 0; i < maxIter; i++ {
		// xy.Mul(x, y)
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
*/

func iterhandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading body")
	}
	lines := strings.Split(string(body), "\n")

	result := ""
	cx := big.NewFloat(0.0)
	cy := big.NewFloat(0.0)
	for idx, line := range lines {
		params := strings.Split(line, ",")
		if len(params) != 1 {

			cx, _, err := cx.Parse(params[0], 10)
			if err != nil {
				log.Error().Err(err).Msg("Error parsing x")
			}
			cy, _, err = cy.Parse(params[1], 10)
			if err != nil {
				log.Error().Err(err).Msg("Error parsing y")
			}
			iter, err := strconv.Atoi(params[2])
			if err != nil {
				log.Error().Err(err).Msg("Error parsing i")
			}

			_, iter = mandelbrot.Mandelbrot(cx, cy, iter)

			if idx == 0 {
				result = fmt.Sprintf("%d", iter)
			} else {
				result = fmt.Sprintf("%s,%d", result, iter)
			}
		}
	}

	fmt.Fprintf(w, "%s", result)

	/*

		cx := big.NewFloat(0.0)
		cy := big.NewFloat(0.0)

		log.Info().Msgf("node(%v) x=%s y=%s i=%s\n", nodeId.String(), r.URL.Query().Get("x"), r.URL.Query().Get("y"), r.URL.Query().Get("i"))

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

		log.Info().Msgf("node(%v) value=%v iter=%d", nodeId.String(), value.String(), iter)
		fmt.Fprintf(w, "%v\n%d\n", value, iter)
	*/
}

func main() {

	log.Info().Msg("Setup Handler")
	http.HandleFunc("/iter", iterhandler)

	log.Info().Msg("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Unexpected error")
}
