package main

import (
	"fmt"
	"math/big"
)

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

func (c Camera) Step() *big.Float {
	return big.NewFloat(0).Quo(c.Scale, big.NewFloat(200))
}

func (c Camera) Right() {
	c.X.Add(c.X, c.Step())
}

func (c Camera) Left() {
	c.X.Sub(c.X, c.Step())
}

func (c Camera) Up() {
	c.Y.Add(c.Y, c.Step())
}

func (c Camera) Down() {
	c.Y.Sub(c.Y, c.Step())
}

func (c Camera) In() {
	fmt.Printf("%d %s\n", c.Scale.Prec(), c.Scale.Text(byte('f'), 120))
	c.Scale.Mul(c.Scale, big.NewFloat(.95))
}

func (c Camera) Out() {
	c.Scale.Mul(c.Scale, big.NewFloat(1.05))
}
