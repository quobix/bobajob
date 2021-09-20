package bobajob

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"math/rand"
	"time"
)

// copied verbatim from https://www.calhoun.io/creating-random-strings-in-go/
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}

func EncodeGobToBase64(val JobEnvelope) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeBase64ToGob(val string) (*JobEnvelope, error) {
	gobBits, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.Write(gobBits)
	enc := gob.NewDecoder(&buf)
	var je JobEnvelope
	err = enc.Decode(&je)
	return &je, nil
}
