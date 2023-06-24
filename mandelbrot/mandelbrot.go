package mandelbrot

import "math/big"

var four = big.NewFloat(float64(4))
var two = big.NewFloat(float64(2))
var half = big.NewFloat(float64(0.5))

func MandelIteration(cx, cy *big.Float, maxIter int) (*big.Float, int) {

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
