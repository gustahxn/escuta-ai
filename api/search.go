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
	query := r.URL.Query().Get("q")
	
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
}