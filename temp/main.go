package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
)

const (
	width, height = 1024, 1024 // Image dimensions
	xmin, xmax    = -2.5, 1    // X-axis range
	ymin, ymax    = -1.5, 1.5  // Y-axis range
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		fmt.Println(py)
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin

			i := uint8(isInMandelbrotSet(x, y))

			img.Set(px, py, mandelbrotColor(i))
		}
	}

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func mandelbrotColor(iterations uint8) color.Color {
	const contrast = 15
	return color.Gray{255 - contrast*iterations}
}

const (
	maxIter             = 256                         // Maximum number of iterations
	precisionBits       = 10                          // Number of fixed-point bits
	fixedPointOne       = 1 << precisionBits          // Fixed-point representation of 1
	escapeRadius        = 2 << (2 * precisionBits)    // Escape radius
	escapeRadiusSquared = escapeRadius * escapeRadius // Escape radius squared
)

/*
	func isInMandelbrotSet(x, y float64) int {
		xInt := int(x * float64(fixedPointOne))
		yInt := int(y * float64(fixedPointOne))
		x0, y0 := xInt, yInt

		var zx, zy int
		for i := 0; i < maxIter; i++ {
			zx, zy = fixedPointMultiply(zx, zx)-fixedPointMultiply(zy, zy)+x0, 2*fixedPointMultiply(zx, zy)+y0
			if fixedPointMagnitudeSquared(zx, zy) > escapeRadiusSquared {
				return i
			}
		}

		return maxIter
	}

	func fixedPointMultiply(a, b int) int {
		return (a * b) >> precisionBits
	}

	func fixedPointMagnitudeSquared(x, y int) int {
		return fixedPointMultiply(x, x) + fixedPointMultiply(y, y)
	}
*/

var ers = new(big.Int).SetInt64(escapeRadiusSquared)

func isInMandelbrotSet(x, y float64) int {
	xInt := int64(x * float64(fixedPointOne))
	yInt := int64(y * float64(fixedPointOne))
	x0 := new(big.Int).SetInt64(xInt)
	y0 := new(big.Int).SetInt64(yInt)

	zx := new(big.Int)
	zy := new(big.Int)
	for i := 0; i < maxIter; i++ {

		//zx = fixedPointMultiply(zx, zx) - fixedPointMultiply(zy, zy) + x0
		zx = new(big.Int).Sub(fixedPointMultiply(zx, zx), fixedPointMultiply(zy, zy))
		zx = zx.Add(zx, x0)

		//zy = 2*fixedPointMultiply(zx, zy) + y0
		zy = new(big.Int).Mul(fixedPointMultiply(zx, zy), new(big.Int).SetInt64(2))
		zy = zy.Add(zy, y0)

		if fixedPointMagnitudeSquared(zx, zy).Cmp(ers) > 0 {
			return i
		}
	}

	return maxIter
}

func fixedPointMultiply(a, b *big.Int) *big.Int {
	temp := new(big.Int).Mul(a, b)
	temp = temp.Rsh(temp, precisionBits)

	return temp
}

func fixedPointMagnitudeSquared(x, y *big.Int) *big.Int {
	x2 := fixedPointMultiply(x, x)
	y2 := fixedPointMultiply(y, y)
	return x2.Add(x2, y2)
}
