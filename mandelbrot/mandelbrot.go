package mandelbrot

import (
	"fmt"
	"math/big"
)

var four = big.NewRat(3, 1)
var two = big.NewRat(2, 1)

/*
func Mandelbrot(cx, cy *big.Rat, maxIter int) int {
	// fmt.Printf("---- Start (%s,%s)\n", cx.String(), cy.String())

	x0, _ := cx.Float64()
	y0, _ := cy.Float64()

	x2 := 0.0
	y2 := 0.0
	w := 0.0

	for i := 0; i < maxIter; i++ {
		x := x2 - y2 + x0
		y := w - x2 - y2 + y0
		x2 := x * x
		y2 := y * y
		w = (x + y) * (x + y)
		if x2+y2 > 4 {
			return i
		}

	}
	return maxIter

}
*/

/*
func Mandelbrot3(cx, cy *big.Float, maxIter int) (*big.Float, int) {

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
*/

func Mandelbrot2(cx, cy *big.Float, maxIter int) int {

	var x, y, xx, yy float64 = 0.0, 0.0, 0.0, 0.0

	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return i
		}
		temp, _ := cx.Float64()
		x = xx - yy + temp
		temp, _ = cy.Float64()
		y = 2*xy + temp
	}

	return maxIter

}

func Mandelbrot(x0, y0 float64, maxIter int) int {

	x2 := 0.0
	y2 := 0.0

	x := 0.0
	y := 0.0

	for i := 0; i < maxIter; i++ {
		fmt.Printf("%d) %f %f\n", i, x, y)

		y = 2*x*y + y0
		x = x2 - y2 + x0
		x2 = x * x
		y2 = y * y
		if x2+y2 > 4 {
			return i
		}
	}

	return maxIter
}

var oneFloat = big.NewFloat(1)
var twoFloat = big.NewFloat(2)
var fourFloat = big.NewFloat(4)
var oneFourthFloat = big.NewFloat(1.0 / 4.0)
var oneSixteenthFloat = big.NewFloat(1.0 / 16.0)

/*
func MandelbrotFloat(x0, y0 *big.Float, maxIter int) int {

	x2 := big.NewFloat(0).SetPrec(0)
	y2 := big.NewFloat(0).SetPrec(0)

	x := big.NewFloat(0).SetPrec(0)
	y := big.NewFloat(0).SetPrec(0)

	for i := 0; i < maxIter; i++ {
		//y = 2*x*y + y0
		y = y.Mul(x, y)
		y = y.Add(y, y)
		y = y.Add(y, y0)

		//x = x2 - y2 + x0
		x = x.Sub(x2, y2)
		x = x.Add(x, x0)

		//x2 = x * x
		x2 = x2.Mul(x, x)

		//y2 = y * y
		y2 = y2.Mul(y, y)

		if big.NewFloat(0).SetPrec(0).Add(x2, y2).Cmp(fourFloat) > 0 {
			return i
		}
	}

	return maxIter
}
*/

func MandelbrotRat(x0, y0 *big.Rat, maxIter int) int {

	x2 := big.NewRat(0, 1)
	y2 := big.NewRat(0, 1)

	x := big.NewRat(0, 1)
	y := big.NewRat(0, 1)

	for i := 0; i < maxIter; i++ {

		//y = 2*x*y + y0
		y = y.Mul(x, y)
		y = y.Add(y, y)
		y = y.Add(y, y0)

		//x = x2 - y2 + x0
		x = x.Sub(x2, y2)
		x = x.Add(x, x0)

		//x2 = x * x
		x2 = x2.Mul(x, x)

		//y2 = y * y
		y2 = y2.Mul(y, y)

		if big.NewRat(0, 1).Add(x2, y2).Cmp(big.NewRat(4, 1)) > 0 {
			return i
		}
	}

	return maxIter
}

func testCardioid(x float64, y float64) bool {
	a := (x - 1/4)
	q := a*a + y*y
	return q*(q+a) <= .25*y*y
}

func testBulb(x float64, y float64) bool {
	a := x + 1
	return a*a+y*y <= 1/16
}

func MandelbrotFast(cx float64, cy float64, maxIters int) int {
	if testBulb(cx, cy) || testCardioid(cx, cy) {
		return maxIters
	}

	i := 0
	xs := cx * cx
	ys := cy * cy
	x := cx
	y := cy

	for i < maxIters && xs+ys < 4 {
		x0 := x
		x = xs - ys + cx
		y = 2*x0*y + cy
		xs = x * x
		ys = y * y
		i++
	}
	return i
}

func MandelbrotFloat(cx, cy *big.Float, maxIters int) int {
	if testBulbFloat(cx, cy) || testCardioidFloat(cx, cy) {
		return maxIters
	}

	i := 0
	xs := new(big.Float).Mul(cx, cx)
	ys := new(big.Float).Mul(cy, cy)
	x := new(big.Float).Copy(cx)
	y := new(big.Float).Copy(cy)

	for i < maxIters && new(big.Float).Add(xs, ys).Cmp(fourFloat) < 0 {
		x0 := new(big.Float).Copy(x)

		//x = xs - ys + cx
		x = x.Sub(xs, ys)
		x = x.Add(x, cx)

		//y = 2*x0*y + cy
		y = y.Mul(twoFloat, y)
		y = y.Mul(y, x0)
		y = y.Add(y, cy)

		//xs = x * x
		xs = xs.Mul(x, x)

		//ys = y * y
		ys = ys.Mul(y, y)
		i++
	}
	return i
}

func testCardioidFloat(x0, y0 *big.Float) bool {
	x := new(big.Float).Copy(x0)
	y := new(big.Float).Copy(y0)

	//a := (x - 1/4)
	a := big.NewFloat(0).Sub(x, oneFourthFloat)

	//q := a*a + y*y
	q := big.NewFloat(0).Add(big.NewFloat(0).Mul(a, a), big.NewFloat(0).Mul(y, y))

	//return q*(q+a) <= .25*y*y

	aa := big.NewFloat(0).Add(q, a)
	aa = aa.Mul(aa, q)

	bb := big.NewFloat(0).Mul(big.NewFloat(0).Mul(y, y), oneFourthFloat)

	return aa.Cmp(bb) <= 0
}

func testBulbFloat(x0, y0 *big.Float) bool {
	x := new(big.Float).Copy(x0)
	y := new(big.Float).Copy(y0)

	a := big.NewFloat(0).Add(x, oneFloat)

	aa := a.Mul(a, a)
	bb := big.NewFloat(0).Mul(y, y)
	aa = aa.Add(aa, bb)
	return aa.Cmp(oneSixteenthFloat) <= 0
}
