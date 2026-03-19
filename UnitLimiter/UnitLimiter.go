package UnitLimiter

import (
	"sync/atomic"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Keeps unit counter and allows units to wait until said counter drops below the limit
type Limiter struct {
	currentCounter int64
	limit          int64
	description    string
}

func MakeUnitLimiter(limit int64, description string) Limiter {
	return Limiter{
		currentCounter: 0,
		limit:          limit,
		description:    description,
	}
}

func (l *Limiter) Acquire() {
	for atomic.LoadInt64(&l.currentCounter) >= l.limit {
		time.Sleep(100 * time.Millisecond)
	}
	atomic.AddInt64(&l.currentCounter, 1)
	// log.Printf("Acquired %s unit %d\n", l.description, atomic.LoadInt64(&l.currentCounter))
}

func (l *Limiter) Release() {
	// log.Printf("Released %s unit %d\n", l.description, atomic.LoadInt64(&l.currentCounter))
	atomic.AddInt64(&l.currentCounter, -1)
}
