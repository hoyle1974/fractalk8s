package main

import (
	"math/big"
	"sync"
)

type Camera struct {
	sync.Mutex
	X, Y, Scale *big.Float
	dirty       bool
}

func NewCamera(x, y, scale *big.Float) *Camera {
	return &Camera{
		X:     x,
		Y:     y,
		Scale: scale,
		dirty: true,
	}
}

func (c *Camera) IsDirty() bool {
	c.Lock()
	defer c.Unlock()

	v := c.dirty
	c.dirty = false
	return v
}

func (c *Camera) step() *big.Float {
	c.dirty = true
	return big.NewFloat(0).Quo(c.Scale, big.NewFloat(200))
}

func (c *Camera) Right() {
	c.Lock()
	defer c.Unlock()
	c.dirty = true
	c.X.Add(c.X, c.step())
}

func (c *Camera) Left() {
	c.Lock()
	defer c.Unlock()

	c.dirty = true
	c.X.Sub(c.X, c.step())
}

func (c *Camera) Up() {
	c.Lock()
	defer c.Unlock()

	c.dirty = true
	c.Y.Sub(c.Y, c.step())
}

func (c *Camera) Down() {
	c.Lock()
	defer c.Unlock()

	c.dirty = true
	c.Y.Add(c.Y, c.step())
}

func (c *Camera) In() {
	c.Lock()
	defer c.Unlock()

	c.dirty = true
	c.Scale.Mul(c.Scale, big.NewFloat(.95))
}

func (c *Camera) Out() {
	c.Lock()
	defer c.Unlock()

	c.dirty = true
	c.Scale.Mul(c.Scale, big.NewFloat(1.05))
}
