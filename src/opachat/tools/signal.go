package tools

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
)

const compress = false

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)

	if err != nil {
		Danger("unzip error", err)
	}

	r, err := gzip.NewReader(&b)

	if err != nil {
		Danger("unzip new reader", err)
	}

	res, err := io.ReadAll(r)

	if err != nil {
		Danger("unzip read all", err)
	}

	return res
}

func zip(in []byte) []byte {
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)

	if err != nil {
		Danger("zip write", err)
	}

	err = gz.Flush()

	if err != nil {
		Danger("zip flush", err)
	}

	err = gz.Close()

	if err != nil {
		Danger("zip close", err)
	}

	return b.Bytes()
}

// DecodeSdp decodes the input from base64
// It can optionally unzip the input after decoding
func DecodeSdp(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)

	if err != nil {
		Danger("decode sdp", err)
	}

	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)

	if err != nil {
		Danger("decode unmarshal sdp", err)
	}
}

// EncodeSdp encodes the input in base64
// It can optionally zip the input before encoding
func EncodeSdp(obj interface{}) string {
	b, err := json.Marshal(obj)

	if err != nil {
		Danger("encode marshal", err)
	}

	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}
