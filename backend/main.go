package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	apiKey := os.Getenv("LASTFM_API_KEY")

	if apiKey == "" {
		log.Println("AVISO: LASTFM_API_KEY não encontrada no .env")
	}

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		log.Printf("Buscando: %s", query)

		apiURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.search&track=%s&api_key=%s&format=json&limit=10", url.QueryEscape(query), apiKey)
		
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tracks)
	})

	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		artist := r.URL.Query().Get("artist")
		track := r.URL.Query().Get("track")
		log.Printf("Similares para: %s - %s", artist, track)

		apiURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.getsimilar&artist=%s&track=%s&api_key=%s&format=json&limit=30", url.QueryEscape(artist), url.QueryEscape(track), apiKey)
		
		resp, _ := http.Get(apiURL)
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		similars, _ := result["similartracks"].(map[string]interface{})
		tracks := similars["track"]

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tracks)
	})

	handler := cors.AllowAll().Handler(http.DefaultServeMux)

	fmt.Println("Backend 'Escuta Aí' rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}