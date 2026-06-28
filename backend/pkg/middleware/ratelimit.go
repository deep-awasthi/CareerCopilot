package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu       sync.Mutex
	clients  map[string]*rateLimiter
	r        rate.Limit
	b        int
	cleanTTL time.Duration
}

func NewRateLimiterStore(r rate.Limit, b int) *RateLimiterStore {
	store := &RateLimiterStore{
		clients:  make(map[string]*rateLimiter),
		r:        r,
		b:        b,
		cleanTTL: 3 * time.Minute,
	}
	go store.cleanup()
	return store
}

func (s *RateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(s.r, s.b)
		s.clients[ip] = &rateLimiter{limiter, time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (s *RateLimiterStore) cleanup() {
	for {
		time.Sleep(time.Minute)
		s.mu.Lock()
		for ip, v := range s.clients {
			if time.Since(v.lastSeen) > s.cleanTTL {
				delete(s.clients, ip)
			}
		}
		s.mu.Unlock()
	}
}

func RateLimiter(requests int, duration time.Duration) gin.HandlerFunc {
	if requests == 0 {
		requests = 100
	}
	if duration == 0 {
		duration = time.Minute
	}
	r := rate.Limit(float64(requests) / duration.Seconds())
	store := NewRateLimiterStore(r, requests)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := store.getLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "too many requests, please slow down",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
