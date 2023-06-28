package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/hoyle1974/fractalk8s/mandelbrot"
	"github.com/rs/zerolog/log"
)

func iterhandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Version 1")

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

			iter = mandelbrot.MandelbrotFloat(cx, cy, iter)

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
