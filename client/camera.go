package main

import "math/big"

type Camera struct {
	X, Y, Scale *big.Float
}

func NewCamera(x, y, scale *big.Float) Camera {
	return Camera{
		X:     x,
		Y:     y,
		Scale: scale,
	}
}
