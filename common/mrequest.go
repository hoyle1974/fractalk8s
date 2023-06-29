package common

import (
	"math/big"
)

type MRequest struct {
	X, Y, Chunk, ScreenWidth, ScreenHeight, Iter int
	CenterX, CenterY, Size                       string
}

func NewMRequest(x, y, chunk, screenWidth, screenHeight int, centerX, centerY, size *big.Float, iter int) MRequest {

	return MRequest{
		X:            x,
		Y:            y,
		Chunk:        chunk,
		ScreenWidth:  screenHeight,
		ScreenHeight: screenHeight,
		Iter:         iter,
		CenterX:      centerX.String(),
		CenterY:      centerY.String(),
		Size:         size.String(),
	}
}

func (m MRequest) ExtractFloats() (*big.Float, *big.Float, *big.Float) {
	cx, _, _ := new(big.Float).Parse(m.CenterX, 10)
	cy, _, _ := new(big.Float).Parse(m.CenterY, 10)
	s, _, _ := new(big.Float).Parse(m.Size, 10)

	return cx, cy, s
}

// func (m MRequest) ToJsonString() string {
// 	jsonString, _ := json.Marshal(m)
// 	return string(jsonString)
// }

// func NewMRequestFromJson(jsonString string) MRequest {
// 	m := MRequest{}
// 	err := json.Unmarshal([]byte(jsonString), &m)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return m
// }

func (m MRequest) ToBytes() []byte {
	return writeAsBytes(m)
}

func NewMRequestFromBytes(b []byte) MRequest {
	m := MRequest{}
	readAsBytes(b, &m)
	return m
}
