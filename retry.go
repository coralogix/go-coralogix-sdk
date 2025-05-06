package coralogix

import (
	"time"

	"github.com/jpillora/backoff"
)

type Retryer interface {
	// Attempt returns the current attempt number.
	Attempt() uint

	// MaxAttempts returns the maximum number of attempts that can be made for
	// an attempt before failing. A value of 0 implies that the attempt should
	// be retried until it succeeds if the errors are retryable.
	MaxAttempts() uint

	// RetryDelay calculates the delay before the next retry attempt.
	// It returns the duration of the delay and a boolean indicating whether the retry should proceed.
	RetryDelay() (time.Duration, bool)
}

type BackoffRetryer struct {
	b           backoff.Backoff
	maxAttempts uint
}

func (br *BackoffRetryer) Attempt() uint {
	return uint(br.b.Attempt())
}

func (br *BackoffRetryer) MaxAttempts() uint {
	return br.maxAttempts
}

func (br *BackoffRetryer) RetryDelay() (time.Duration, bool) {
	if br.maxAttempts > 0 && uint(br.b.Attempt()) >= br.maxAttempts {
		return 0, false
	}

	return br.b.Duration(), true
}

type InfiniteBackoffRetryer struct {
	b backoff.Backoff
}

func (ibr *InfiniteBackoffRetryer) Attempt() uint {
	return uint(ibr.b.Attempt())
}

func (ibr *InfiniteBackoffRetryer) MaxAttempts() uint {
	return 0
}

func (ibr *InfiniteBackoffRetryer) RetryDelay() (time.Duration, bool) {
	return ibr.b.Duration(), true
}
