package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// song structure
type Song struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
}

// Global variables
var songs []Song
var nextID = 1

func main() {
	// register handlers for song request
	http.HandleFunc("/songs", handleSongs)
	http.HandleFunc("/songs/", handleSongByID)

	fmt.Println("Server is running on port 8579")
	log.Fatal(http.ListenAndServe(":8579", nil))
}

// handleSongs handles request to the endpoint "/songs"
func handleSongs(w http.ResponseWriter, r *http.Request) {
	// Check for GET or POST methods
	if r.Method == http.MethodGet {
		getAllSongs(w, r)
	} else if r.Method == http.MethodPost {
		createSong(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSongByID handles requests to the endpoint "/songs/{id}"
func handleSongByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path and convert to integer
	path := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Handle requests based on methods
	switch r.Method {
	case http.MethodGet:
		getSongByID(w, r, id)
	case http.MethodPut:
		updateSong(w, r, id)
	case http.MethodDelete:
		deleteSong(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllSongs retrieves all songs and sends them as JSON response
func getAllSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// getSongByID retrieves a specific song by ID and sends it as JSON response
func getSongByID(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	for _, song := range songs {
		if song.ID == id {
			json.NewEncoder(w).Encode(song)
			return
		}
	}
	http.Error(w, "Song not found", http.StatusNotFound)
}

// createSong creates a new song from request body and adds it to the list
func createSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var song Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	song.ID = nextID
	nextID++

	songs = append(songs, song)
	json.NewEncoder(w).Encode(song)
}

// updateSong updates an existing song by ID based on request body
func updateSong(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	for i, song := range songs {
		if song.ID == id {
			var updatedSong Song
			err := json.NewDecoder(r.Body).Decode(&updatedSong)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			updatedSong.ID = id // Keep the same ID
			songs[i] = updatedSong
			json.NewEncoder(w).Encode(updatedSong)
			// updateSong(w, r, id)
			json.NewEncoder(w).Encode(updatedSong)
			return
		}
	}
	http.Error(w, "Song not found", http.StatusNotFound)
}

// deleteSong deletes an existing song by ID
func deleteSong(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	for i, song := range songs {
		if song.ID == id {
			songs = append(songs[:i], songs[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{"message": "Song deleted"})
			return
		}
	}
	http.Error(w, "Song not found", http.StatusNotFound)
}
