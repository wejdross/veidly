package helpers

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
)

/*
	generate cRNG number x such as maxExlusive > x >= 0
*/
func CRNG_uint32(maxExclusive uint32) (uint32, error) {
	buf := make([]byte, 4)
	var err error

	/* TODO: loop it again on error. Dont give up so easily */
	if _, err = rand.Read(buf); err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint32(buf) % maxExclusive, err
}

/*
	same as RngCryptoInt but will panic on error
*/
func CRNG_uint32Panic(maxExclusive uint32) uint32 {
	ret, err := CRNG_uint32(maxExclusive)
	if err != nil {
		panic(err)
	}
	return ret
}

/*
	time.sleep for random time t such as maxMilisec > t >= 0
	on error do nothing
*/
func CRNG_Delay(maxMilisec uint32) {
	if sleep, tmperr := CRNG_uint32(maxMilisec); tmperr == nil {
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
}

/*
	flip random byte in string
*/
func CRNG_ReplaceByte(v string) string {
	if len(v) == 0 {
		return ""
	}
	c := []byte(v)
	ix := CRNG_uint32Panic(uint32(len(c)))
	max := uint32(len(__CRNG_Charset))
	c[ix] = __CRNG_Charset[CRNG_uint32Panic(max)]
	return string(c)
}

/*
	random []byte
*/
func CRNG_Bytes(len int) ([]byte, error) {
	buf := make([]byte, len)
	var err error

	/* TODO: loop it again on error. Dont give up so easily */
	if _, err = rand.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

/*
	random []byte or panic
*/
func CRNG_BytesPanic(len int) []byte {
	r, err := CRNG_Bytes(len)
	if err != nil {
		panic(err)
	}
	return r
}

const __CRNG_Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func CRNG_string(n int64) (string, error) {
	b := make([]byte, n)
	max := uint32(len(__CRNG_Charset))
	for i := range b {
		d, err := CRNG_uint32(max)
		if err != nil {
			return "", err
		}
		b[i] = __CRNG_Charset[d]
	}
	return string(b), nil
}

/*
	Generate randopm n characters
*/
func CRNG_stringPanic(n int64) string {
	s, err := CRNG_string(n)
	if err != nil {
		panic(err)
	}
	return s
}

// generate safe and unique sequence of characters
func GetUniqueToken() string {
	return uuid.New().String()
}
