package common

import (
	"time"
)

type MResponse struct {
	Iter     []int
	CalcTime time.Duration
}

func NewMResponse(iter []int, calcTime time.Duration) MResponse {

	m := MResponse{
		Iter:     iter,
		CalcTime: calcTime,
	}

	return m
}

func (m MResponse) Extract() ([]int, time.Duration) {
	return m.Iter, m.CalcTime
}

func (m MResponse) ToBytes() []byte {
	return writeAsBytes(m, compressResponse)
}

func NewMResponseFromBytes(b []byte) MResponse {
	m := MResponse{}
	readAsBytes(b, &m, compressResponse)
	return m
}
