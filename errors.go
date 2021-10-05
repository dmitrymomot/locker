package locker

import "errors"

// Predefined package errors
var (
	ErrAlreadyTaken     = errors.New("the lock key is already in use")
	ErrAlreadyUnlocked  = errors.New("the key was already expired or unlocked")
	ErrTryUnlockAnother = errors.New("could not unlock key which was locked by another process")
	ErrExpired          = errors.New("lock is expired")
)
