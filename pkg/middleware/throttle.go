package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// Throttle holds relevant parameters for configuring how the throttle middleware behaves
type Throttle struct {
	IPs     map[string]ThrottleClient
	Mu      *sync.RWMutex
	Rate    rate.Limit
	Burst   int
	Cleanup time.Duration
}

// ThrottleClient holds info about an ip such as it's rate limiter and the last time a request was made
type ThrottleClient struct {
	Limiter *rate.Limiter
	last    time.Time
}

func (t Throttle) add(ip string) *rate.Limiter {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	limiter := rate.NewLimiter(t.Rate, t.Burst)
	t.IPs[ip] = ThrottleClient{
		Limiter: limiter,
		last:    time.Now(),
	}
	return limiter
}

func (t Throttle) getLimiter(ip string) *rate.Limiter {
	t.Mu.Lock()
	c, exists := t.IPs[ip]
	if !exists {
		t.Mu.Unlock()
		return t.add(ip)
	}
	c.last = time.Now()
	t.IPs[ip] = c
	t.Mu.Unlock()
	return c.Limiter
}

func (t Throttle) cleanup() {
	t.Mu.Lock()
	for ip, v := range t.IPs {
		if time.Since(v.last) > t.Cleanup {
			delete(t.IPs, ip)
		}
	}
	t.Mu.Unlock()
}

// Throttle prevents multiple repeated requests in a certain time period
func (t Throttle) Throttle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO if we are behind a load balancer we need to support X-Forwarded-For headers
	limiter := t.getLimiter(r.RemoteAddr)
	if !limiter.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	next.ServeHTTP(w, r)
}
