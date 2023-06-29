package main

import (
	"fmt"
	"io/ioutil"
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

	request := common.NewMRequestFromJson(string(body))
	cx, cy, iter := request.Extract()

	resultIter := make([]int, len(cx))

	for idx, _ := range cx {
		resultIter[idx] = mandelbrot.MandelbrotFloat(cx[idx], cy[idx], iter)
	}

	end := time.Now()
	calcTime := end.Sub(start)
	response := common.NewMResponse(resultIter, calcTime)

	fmt.Fprintf(w, "%s", response.ToJsonString())
}

func main() {

	log.Info().Msg("Setup Handler")
	http.HandleFunc("/iter", iterhandler)

	log.Info().Msg("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Error().Err(err).Msg("Unexpected error")
}
