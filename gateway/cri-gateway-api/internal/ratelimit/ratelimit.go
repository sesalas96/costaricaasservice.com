// Package ratelimit implementa un rate limiter sliding-window simple en
// memoria, por clave (IP o citizen sub).
//
// Diseño: token bucket por clave. Cada bucket almacena tokens de capacidad
// `burst` y se rellena a `rate` tokens por segundo. Limpiezas periódicas
// quitan buckets sin actividad para evitar leak de memoria.
package ratelimit

import (
	"net/http"
	"sync"
	"time"

	scrctx "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/ctx"
	screrrors "github.com/devsebas/costaricaasservice/libs/cri-lib-shared/errors"
	"github.com/devsebas/costaricaasservice/libs/cri-lib-shared/httpx"

	authmw "github.com/devsebas/costaricaasservice/libs/cri-lib-auth/middleware"
)

// Limiter es un token-bucket rate limiter.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64       // tokens/segundo
	burst   float64       // capacidad
	idleTTL time.Duration // limpieza de buckets sin actividad
}

type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// New construye un Limiter.
func New(rate, burst float64, idleTTL time.Duration) *Limiter {
	if idleTTL == 0 {
		idleTTL = 5 * time.Minute
	}
	l := &Limiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		burst:   burst,
		idleTTL: idleTTL,
	}
	go l.gcLoop()
	return l
}

// Allow consume un token de la bucket de `key`. Retorna false si no hay
// suficientes (rechazar request).
func (l *Limiter) Allow(key string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &bucket{tokens: l.burst - 1, lastSeen: now}
		return true
	}
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.lastSeen = now
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

func (l *Limiter) gcLoop() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for range t.C {
		cutoff := time.Now().Add(-l.idleTTL)
		l.mu.Lock()
		for k, b := range l.buckets {
			if b.lastSeen.Before(cutoff) {
				delete(l.buckets, k)
			}
		}
		l.mu.Unlock()
	}
}

// Middleware enforce el rate limit. La key es citizen sub (si autenticado)
// o IP del cliente.
func Middleware(l *Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get(authmw.HeaderSub)
			if key == "" {
				key = clientIP(r)
			}
			if !l.Allow(key) {
				httpx.Fail(w, r, screrrors.New(screrrors.CodeRateLimited, "rate limit exceeded"))
				return
			}
			_ = scrctx.RequestID(r.Context()) // marca uso (evita unused import warning si no se usa)
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := indexComma(xff); i > 0 {
			return xff[:i]
		}
		return xff
	}
	if r.RemoteAddr != "" {
		if i := indexColon(r.RemoteAddr); i > 0 {
			return r.RemoteAddr[:i]
		}
		return r.RemoteAddr
	}
	return "unknown"
}

func indexComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}

func indexColon(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			return i
		}
	}
	return -1
}
