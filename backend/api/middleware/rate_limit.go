package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// RateStore represents a Redis-compatible backend for distributed rate limiting.
type RateStore interface {
	Allow(key string, limit rate.Limit, burst int) bool
}

// MemoryRateStore provides an in-memory fallback mimicking Redis.
type MemoryRateStore struct {
	mu      sync.Mutex
	clients map[string]*client
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewMemoryRateStore() *MemoryRateStore {
	store := &MemoryRateStore{clients: make(map[string]*client)}
	go func() {
		for {
			time.Sleep(time.Minute)
			store.mu.Lock()
			for k, c := range store.clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(store.clients, k)
				}
			}
			store.mu.Unlock()
		}
	}()
	return store
}

func (s *MemoryRateStore) Allow(key string, limit rate.Limit, burst int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, found := s.clients[key]; !found {
		s.clients[key] = &client{limiter: rate.NewLimiter(limit, burst)}
	}
	s.clients[key].lastSeen = time.Now()
	return s.clients[key].limiter.Allow()
}

// Global default store
var defaultStore RateStore = NewMemoryRateStore()

// RateLimit implements key-based rate limiting (user_id:tenant_id or IP fallback).
func RateLimit(reqsPerSec float64, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr // default IP

			// Attempt to extract composite key if authenticated
			userVal := r.Context().Value(UserIDKey)
			tenantIDStr := GetTenantID(r.Context())
			
			if userVal != nil && tenantIDStr != "" {
				if userID, ok := userVal.(uuid.UUID); ok {
					key = tenantIDStr + ":" + userID.String()
				}
			}

			if !defaultStore.Allow(key, rate.Limit(reqsPerSec), burst) {
				slog.Warn("Rate limit exceeded", "key", key)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "rate_limit_exceeded", "message": "Too many requests"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}


