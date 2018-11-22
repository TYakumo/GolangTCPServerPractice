package tcpserver

import (
	"sync"
	"time"
)

type RateLimitController struct {
	availableToken int
	tokenLimit     int
	lastReplTime   time.Time
	replenishRate  time.Duration
	m              sync.Mutex
}

func (r *RateLimitController) ReplenishToken() {
	timeDiff := time.Now().Sub(r.lastReplTime)

	if timeDiff >= r.replenishRate {
		numberOfToken := int(timeDiff / r.replenishRate)
		r.availableToken = r.availableToken + numberOfToken
		if r.availableToken > r.tokenLimit {
			r.availableToken = r.tokenLimit
		}
		r.lastReplTime = time.Now()
	}
}

func (r *RateLimitController) GetToken() bool {
	r.m.Lock()
	defer r.m.Unlock()

	r.ReplenishToken()
	if r.availableToken > 0 {
		r.availableToken--
		return true
	} else {
		return false
	}
}

func NewRateLimitController(tokenLimit int, replenishRate time.Duration) (*RateLimitController, error) {
	var r RateLimitController
	r.availableToken = tokenLimit
	r.tokenLimit = tokenLimit
	r.lastReplTime = time.Now()
	r.replenishRate = replenishRate
	return &r, nil
}
