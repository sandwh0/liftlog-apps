package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// WorkoutLog represents the expected input payload
type WorkoutLog struct {
	Exercise string  `json:"exercise"`
	Reps     int     `json:"reps"`
	Weight   float64 `json:"weight"`
}

// XPResponse represents the JSON response
type XPResponse struct {
	Exercise string  `json:"exercise"`
	Reps     int     `json:"reps"`
	Weight   float64 `json:"weight"`
	XPGained int     `json:"xp_gained"`
	Timestamp string `json:"timestamp"`
}

// logWorkoutHandler handles POST /log requests
func logWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var workoutLog WorkoutLog
	if err := json.NewDecoder(r.Body).Decode(&workoutLog); err != nil {
		log.Printf("Failed to decode request: %v", err)
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Basic validation
	if workoutLog.Exercise == "" {
		http.Error(w, "Exercise name is required", http.StatusBadRequest)
		return
	}
	if workoutLog.Reps <= 0 {
		http.Error(w, "Reps must be greater than 0", http.StatusBadRequest)
		return
	}
	if workoutLog.Weight <= 0 {
		http.Error(w, "Weight must be greater than 0", http.StatusBadRequest)
		return
	}

	// Calculate XP based on exercise difficulty and volume
	// More complex formula that considers both weight and reps
	baseXP := float64(workoutLog.Reps) * workoutLog.Weight
	volumeMultiplier := 1.0
	if workoutLog.Reps > 12 {
		volumeMultiplier = 0.8 // Lower multiplier for high rep sets
	} else if workoutLog.Reps < 5 {
		volumeMultiplier = 1.2 // Higher multiplier for low rep sets
	}
	xp := int(baseXP * volumeMultiplier * 0.1)

	resp := XPResponse{
		Exercise: workoutLog.Exercise,
		Reps:     workoutLog.Reps,
		Weight:   workoutLog.Weight,
		XPGained: xp,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Logged workout: %s - %d reps @ %.1f kg = %d XP", 
		workoutLog.Exercise, workoutLog.Reps, workoutLog.Weight, xp)
}

func main() {
	// Set up logging with timestamp
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/log", logWorkoutHandler)
	
	port := ":8080"
	fmt.Printf("Starting workout API on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
