package locker

import "time"

// WithTries option sets custom lock tries
func WithTries(t uint) Option {
	if t < 1 {
		t = 1
	}
	return func(m *Mutex) {
		m.tries = t
	}
}

// WithValueGenerationFunc option sets custom value generation function
func WithValueGenerationFunc(fn func() []byte) Option {
	return func(m *Mutex) {
		m.genValueFunc = fn
	}
}

// WithDalayFunc option sets custom retry dalay function
func WithDalayFunc(fn func(tries int) time.Duration) Option {
	return func(m *Mutex) {
		m.delayFunc = fn
	}
}
