package locker

import (
	"math/rand"
	"time"
)

// GenEmptyValue to use when you're sure in your code
// and wouldn't like to store some value into storage
func GenEmptyValue() []byte {
	return []byte("")
}

// generates random value to avoid unlock key,
// which is taken by another process
func defaultGenValue() []byte {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return []byte("")
	}
	return b
}

// generates pseudo-random delay time
// to avoid run all pods at the same time
// after failed first lock attempting
func defaultDelayFunc(tries int) time.Duration {
	if tries < 1 {
		tries = 1
	}
	rn := rand.Intn(tries)
	if rn == 0 {
		rn = 1
	}
	return time.Duration(rn * 100 * int(time.Millisecond))
}
