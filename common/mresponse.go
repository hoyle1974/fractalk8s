package common

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
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
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(m)
	if err != nil {
		panic(err)
	}

	var zipped bytes.Buffer
	gz := gzip.NewWriter(&zipped)
	if _, err := gz.Write(network.Bytes()); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return zipped.Bytes()
}

func NewMResponseFromBytes(b []byte) MResponse {

	gz, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	output, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	network := bytes.NewBuffer(output)

	dec := gob.NewDecoder(network)
	m := MResponse{}
	err = dec.Decode(&m)
	if err != nil {
		panic(err)
	}

	return m
}
