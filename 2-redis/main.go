package main

import (
	"bytes"
	"net/http"
	"time"

	"github.com/ernestii/learning-go/2-redis/cache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	stringCache cache.StringCache
)

func setEndpoint(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	r.Body.Close()
	val := buf.String()
	stringCache.Set(key, val)

	w.Write([]byte("Saved"))
}

func getEndpoint(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	val := stringCache.Get(key)

	if val == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Write([]byte(val))
}

func main() {
	// stringCache = cache.NewLocalCache()
	stringCache = cache.NewRedisCache("localhost:7001", 1, 3*time.Minute)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/set/{key}", setEndpoint)
	r.Get("/get/{key}", getEndpoint)

	http.ListenAndServe(":8081", r)
}
