package common

import (
	"math/big"
)

type MRequest struct {
	X, Y, Chunk, ScreenWidth, ScreenHeight, Iter int
	CenterX, CenterY, Size                       *big.Float
}

func NewMRequest(x, y, chunk, screenWidth, screenHeight int, centerX, centerY, size *big.Float, iter int) MRequest {

	return MRequest{
		X:            x,
		Y:            y,
		Chunk:        chunk,
		ScreenWidth:  screenHeight,
		ScreenHeight: screenHeight,
		Iter:         iter,
		CenterX:      centerX,
		CenterY:      centerY,
		Size:         size,
	}
}

// func (m MRequest) ExtractFloats() (*big.Float, *big.Float, *big.Float) {
// 	cx, _, _ := new(big.Float).Parse(m.CenterX, 10)
// 	cy, _, _ := new(big.Float).Parse(m.CenterY, 10)
// 	s, _, _ := new(big.Float).Parse(m.Size, 10)

// 	return cx, cy, s
// }

func (m MRequest) ToBytes() []byte {
	return writeAsBytes(m, compressRequest)
}

func NewMRequestFromBytes(b []byte) MRequest {
	m := MRequest{}
	readAsBytes(b, &m, compressRequest)
	return m
}
