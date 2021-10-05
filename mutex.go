package locker

import (
	"bytes"
	"fmt"
	"time"
)

type (
	// Mutex / mutual exclusion lock
	Mutex struct {
		storage      mutexStorageDriver
		tries        uint
		genValueFunc func() []byte
		delayFunc    func(tries int) time.Duration
	}

	// a lock value and expiration time
	locked struct {
		key      string
		value    []byte
		until    time.Time
		unlocked bool
	}

	// Option of mutex
	Option func(*Mutex)

	// the interface of a driver to manage lock keys
	mutexStorageDriver interface {
		// Set creates a record in the storage with the key and data and returns
		// numberic ID, true of the created record. If key has already exists,
		// then (0, false, nil) returned.
		Set(key string, data []byte, ttl time.Duration) (int64, bool, error)

		// Get returns the data associated with the key and it's numeric ID
		Get(key string) (int64, []byte, error)

		// Delete deletes a record associated with the given key and returns the numeric ID of the deleted record.
		// numeric ID of the deleted record. If key doesn't exists
		// then (0, nil) returned
		Delete(key string) (int64, error)
	}
)

// New mutex instance with options
func New(s mutexStorageDriver, opt ...Option) *Mutex {
	m := &Mutex{
		storage:      s,
		tries:        1,
		genValueFunc: defaultGenValue,
		delayFunc:    defaultDelayFunc,
	}
	for _, o := range opt {
		o(m)
	}
	if m.storage == nil {
		panic("missed storage driver")
	}
	return m
}

// Lock func
// Returns error if the lock is already in use.
func (m *Mutex) Lock(key string, ttl time.Duration) (func() error, error) {
	val := m.genValueFunc()
	for i := 0; i < int(m.tries); i++ {
		if i != 0 {
			time.Sleep(m.delayFunc(i))
		}

		start := time.Now()

		_, ok, err := m.storage.Set(key, val, ttl)
		if err != nil {
			return nil, fmt.Errorf("could not store key=%s with value=%v: %w", key, val, err)
		}
		if ok {
			now := time.Now()
			until := now.Add(ttl - now.Sub(start))

			if now.Before(until) {
				l := &locked{
					key:   key,
					value: val,
					until: until,
				}
				return func() error {
					return m.unlock(l)
				}, nil
			}

			return nil, ErrExpired
		}
	}

	return nil, ErrAlreadyTaken
}

// Unlock func. Returns error if the lock is
// already in use by another process, or expired.
func (m *Mutex) unlock(l *locked) error {
	if l.unlocked {
		return ErrAlreadyUnlocked
	}

	if time.Now().After(l.until) {
		return ErrExpired
	}

	_, res, err := m.storage.Get(l.key)
	if err != nil {
		return fmt.Errorf("could not get lock record by key=%s: %w", l.key, err)
	}

	if bytes.Equal(l.value, res) {
		if _, err := m.storage.Delete(l.key); err != nil {
			return fmt.Errorf("could not delete lock record by key=%s: %w", l.key, err)
		}
		l.unlocked = true
		return nil
	}

	return ErrTryUnlockAnother
}
