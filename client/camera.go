package main

import "math/big"

type Camera struct {
	X, Y, Scale *big.Float
}

func NewCamera(x, y, scale float64) Camera {
	return Camera{
		X:     big.NewFloat(x),
		Y:     big.NewFloat(y),
		Scale: big.NewFloat(scale),
	}
}
