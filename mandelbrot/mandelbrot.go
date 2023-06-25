package mandelbrot

import "math/big"

var four = big.NewFloat(float64(4))
var two = big.NewFloat(float64(2))
var half = big.NewFloat(float64(0.5))

func Mandelbrot(cx, cy *big.Float, maxIter int) (*big.Float, int) {

	var x, y, xx, yy *big.Float = big.NewFloat(0), big.NewFloat(0), big.NewFloat(0), big.NewFloat(0)

	for i := 0; i < maxIter; i++ {
		xy := big.NewFloat(0).Mul(x, y)
		xx = big.NewFloat(0).Mul(x, x)
		yy = big.NewFloat(0).Mul(y, y)
		if big.NewFloat(0).Add(xx, yy).Cmp(four) > 0 {
			return big.NewFloat(0).Add(xx, yy), i
		}

		//x = xx - yy + temp
		x = big.NewFloat(0).Sub(xx, yy)
		x = x.Add(x, cx)

		//y = 2*xy + temp
		y = y.Mul(two, xy)
		y = y.Add(y, cy)
	}

	// logZn := (x*x + y*y) / 2
	return big.NewFloat(0), maxIter

}

func Mandelbrot2(cx, cy *big.Float, maxIter int) (*big.Float, int) {

	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return big.NewFloat(xx + yy), i
		}
		temp, _ := cx.Float64()
		x = xx - yy + temp
		temp, _ = cy.Float64()
		y = 2*xy + temp
	}

	logZn := (x*x + y*y) / 2
	return big.NewFloat(logZn), maxIter

}
