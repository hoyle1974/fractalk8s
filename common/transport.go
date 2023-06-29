package common

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
)

const (
	compressRequest  = false
	compressResponse = true
)

func writeAsBytes[K any](m K, compress bool) []byte {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(m)
	if err != nil {
		panic(err)
	}

	if !compress {
		return network.Bytes()
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

func readAsBytes[K any](b []byte, m K, compress bool) {
	var output []byte
	if compress {
		gz, err := gzip.NewReader(bytes.NewBuffer(b))
		if err != nil {
			panic(err)
		}
		output, err = ioutil.ReadAll(gz)
		if err != nil {
			panic(err)
		}
	} else {
		output = b
	}

	network := bytes.NewBuffer(output)

	dec := gob.NewDecoder(network)
	err := dec.Decode(&m)
	if err != nil {
		panic(err)
	}

}
