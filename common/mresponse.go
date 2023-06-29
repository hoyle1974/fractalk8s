package common

import (
	"encoding/json"
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

func (m MResponse) ToJsonString() string {
	jsonString, _ := json.Marshal(m)

	return string(jsonString)
}

func NewMResponseFromJson(jsonString string) MResponse {
	m := MResponse{}
	err := json.Unmarshal([]byte(jsonString), &m)
	if err != nil {
		panic(err)
	}

	return m
}
