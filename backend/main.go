package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type CacheItem struct {
	Data      interface{}
	Timestamp time.Time
}

var (
	cache      = make(map[string]CacheItem)
	cacheMutex sync.RWMutex
	ttl        = 30 * time.Minute
)

func getFromCache(key string) (interface{}, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	item, found := cache[key]
	if found && time.Since(item.Timestamp) < ttl {
		return item.Data, true
	}
	return nil, false
}

func setToCache(key string, data interface{}) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache[key] = CacheItem{Data: data, Timestamp: time.Now()}
}

func main() {
    _ = godotenv.Load()
    
    apiKey := os.Getenv("LASTFM_API_KEY")
    if apiKey == "" {
        log.Fatal("LASTFM_API_KEY não configurada")
    }
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query().Get("q")
        if query == "" {
            http.Error(w, "Query parameter 'q' is required", 400)
            return
        }
        
        cacheKey := "search:" + query
        if data, found := getFromCache(cacheKey); found {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(data)
            return
        }
        
        apiURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.search&track=%s&api_key=%s&format=json&limit=10", 
            url.QueryEscape(query), apiKey)
        
        resp, err := http.Get(apiURL)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        defer resp.Body.Close()
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        res, _ := result["results"].(map[string]interface{})
        matches, _ := res["trackmatches"].(map[string]interface{})
        tracks := matches["track"]
        
        setToCache(cacheKey, tracks)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tracks)
    })

    http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
        artist := r.URL.Query().Get("artist")
        track := r.URL.Query().Get("track")
        
        if artist == "" || track == "" {
            http.Error(w, "Parameters 'artist' and 'track' are required", 400)
            return
        }
        
        cacheKey := "rec:" + artist + ":" + track
        if data, found := getFromCache(cacheKey); found {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(data)
            return
        }
        
        apiURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.getsimilar&artist=%s&track=%s&api_key=%s&format=json&limit=30", 
            url.QueryEscape(artist), url.QueryEscape(track), apiKey)
        
        resp, err := http.Get(apiURL)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        defer resp.Body.Close()
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        similars, _ := result["similartracks"].(map[string]interface{})
        tracks := similars["track"]
        
        setToCache(cacheKey, tracks)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tracks)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    c := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:3000",
            "http://localhost:5173",
            "https://*.vercel.app",  
            "https://escuta-ai.vercel.app", 
        },
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        AllowCredentials: false,
    })

    handler := c.Handler(http.DefaultServeMux)
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}