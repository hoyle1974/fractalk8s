package common

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io/ioutil"
)

func writeAsBytes[K any](m K) []byte {
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

func readAsBytes[K any](b []byte, m K) {
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
	err = dec.Decode(&m)
	if err != nil {
		panic(err)
	}

}
