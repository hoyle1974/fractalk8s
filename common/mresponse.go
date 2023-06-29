package common

import (
	"time"
)

type MResponse struct {
	Iter     []byte
	CalcTime time.Duration
}

func NewMResponse(iter []byte, calcTime time.Duration) MResponse {

	m := MResponse{
		Iter:     iter,
		CalcTime: calcTime,
	}

	return m
}

func (m MResponse) Extract() ([]byte, time.Duration) {
	return m.Iter, m.CalcTime
}

func (m MResponse) ToBytes() []byte {
	return writeAsBytes(m)
}

func NewMResponseFromBytes(b []byte) MResponse {
	m := MResponse{}
	readAsBytes(b, &m)
	return m
}
