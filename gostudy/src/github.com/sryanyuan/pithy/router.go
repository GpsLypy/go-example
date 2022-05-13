package pithy

import (
	"net/http"
	"net/http/pprof"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// SwappableHandler is a http.Handler thats allows you to swap mux during working
type SwappableHandler struct {
	lock        sync.RWMutex
	router      *mux.Router
	enablePprof bool
}

// Swap swaps the router of the handler
func (h *SwappableHandler) Swap(r *mux.Router) {
	h.lock.Lock()
	h.router = r
	h.lock.Unlock()
}

// ServeHTTP implement the http.Handler interface
func (h *SwappableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.filter(w, r) {
		return
	}
	h.lock.RLock()
	router := h.router
	h.lock.RUnlock()
	router.ServeHTTP(w, r)
}

func (h *SwappableHandler) filter(w http.ResponseWriter, r *http.Request) bool {
	// pprof support, enable by default
	if h.enablePprof &&
		strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		pprof.Index(w, r)
		return true
	}
	return false
}

// APIRouter defines an interface to determine which handler is called by path and method
type APIRouter interface {
	Path() string
	Method() string
	Handler() APIFunc
	Version() bool
}
