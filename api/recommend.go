package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("LASTFM_API_KEY")
	artist := r.URL.Query().Get("artist")
	track := r.URL.Query().Get("track")

	apiURL := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.getsimilar&artist=%s&track=%s&api_key=%s&format=json&limit=30", url.QueryEscape(artist), url.QueryEscape(track), apiKey)
	
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}