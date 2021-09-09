package patterns

import (
	"log"
	"math/rand"
	"net/http"
)

func CurrentQueueDepth() int {
	return rand.Intn(2000)
}

const MaxQueueDepth = 1000

func LoadShedding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CurrentQueueDepth() > MaxQueueDepth {
			log.Printf("load shedding engaged")
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		next.ServeHTTP(w, r)
	})
}
