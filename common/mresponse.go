package common

import (
	"encoding/json"
)

type MResponse struct {
	Iter []int
}

func NewMResponse(iter []int) MResponse {

	m := MResponse{
		Iter: iter,
	}

	return m
}

func (m MResponse) Extract() []int {
	return m.Iter
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
