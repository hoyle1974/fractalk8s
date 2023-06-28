package common

import (
	"encoding/json"
	"math/big"
)

type MRequest struct {
	X    []string
	Y    []string
	Iter int
}

func NewMRequest(x, y []*big.Float, iter int) MRequest {

	m := MRequest{
		X:    make([]string, len(x)),
		Y:    make([]string, len(y)),
		Iter: iter,
	}

	for i := 0; i < len(x); i++ {
		m.X[i] = x[i].String()
		m.Y[i] = y[i].String()
	}

	return m
}

func (m MRequest) Extract() ([]*big.Float, []*big.Float, int) {
	x := make([]*big.Float, len(m.X))
	y := make([]*big.Float, len(m.Y))

	for i := 0; i < len(x); i++ {
		x[i], _, _ = new(big.Float).Parse(m.X[i], 10)
		y[i], _, _ = new(big.Float).Parse(m.Y[i], 10)
	}

	return x, y, m.Iter

}

func (m MRequest) ToJsonString() string {
	jsonString, _ := json.Marshal(m)

	return string(jsonString)
}

func NewMRequestFromJson(jsonString string) MRequest {
	m := MRequest{}
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		panic(err)
	}

	return m
}
