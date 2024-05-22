package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	locations := make(map[string]string)

	local, _ := time.LoadLocation("Local")
	locations["Local"] = getCurTime(local)
	var err error

	if queries.Has("location") {
		reqLocations := queries["location"]
		for _, reqLocation := range reqLocations {
			var location *time.Location
			if location, err = time.LoadLocation(reqLocation); err != nil {
				locations[reqLocation] = err.Error()
			} else {
				locations[reqLocation] = getCurTime(location)
			}
		}
	}
	// b := new(bytes.Buffer)
	// for key, val := range locations {
	// 	fmt.Fprintf(b, "%s: %v\n", key, val)
	// }
	// fmt.Fprintf(w, b.String())
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getCurTime(location *time.Location) string {
	return time.Now().In(location).String()
}

func main() {
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// This is a comment.
