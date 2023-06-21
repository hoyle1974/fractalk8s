package main

import (
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/alitto/pond"
)

func calc(x, y *big.Float) int {

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/iter?x=%s&y=%s&i=10", x.String(), y.String()))
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	ii := strings.Split(string(body), "\n")[1]
	i, _ := strconv.Atoi(ii)
	return i
}

func main() {

	temp := " .-:[]{}*0@"

	// m := mandelbrot.Mandelbrot(mandelbrot.Min, mandelbrot.Max, 20, 60, 9)

	// for i := range m {
	// 	fmt.Println(m[i])
	// }

	const mx = 150
	const my = 50

	var values [mx][my]int

	pool := pond.New(100, mx*my+1)

	for ty := 0; ty < my; ty++ {
		for tx := 0; tx < mx; tx++ {
			x := tx
			y := ty
			pool.Submit(func() {

				xx := big.NewFloat(float64(x))
				yy := big.NewFloat(float64(y))
				xx.Quo(xx, big.NewFloat(mx))
				xx.Mul(xx, big.NewFloat(2.8))
				xx.Sub(xx, big.NewFloat(2.0))

				yy.Quo(yy, big.NewFloat(my))
				yy.Mul(yy, big.NewFloat(2.8))
				yy.Sub(yy, big.NewFloat(1.4))

				values[x][y] = calc(xx, yy)
			})
		}
	}

	for pool.CompletedTasks() != mx*my {

		for y := 0; y < my; y++ {
			for x := 0; x < mx; x++ {
				fmt.Printf("%c", temp[values[x][y]])
			}
			fmt.Println("")

		}
	}
	pool.StopAndWait()

	/*
		for x := -2.0; x < 0.8; x += 0.075 {
			for y := -1.4; y < 1.4; y += 0.015 {
				resp, err := http.Get(fmt.Sprintf("http://localhost:8080/iter?x=%f&y=%f&i=10", x, y))
				if err != nil {
					panic(err)
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
				ii := strings.Split(string(body), "\n")[1]
				i, _ := strconv.Atoi(ii)
				fmt.Printf("%c", temp[i])
			}
			fmt.Println("")
		}
	*/

}
