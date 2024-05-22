package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type LocationTime struct {
	Location string
	Time     string
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	locations := getLocationsFromQueries(queries)
	locationTimes := buildLocationTimes(locations)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locationTimes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildLocationTimes(locations map[string]string) []*LocationTime{
	var locationTimes []*LocationTime
	for key, val := range locations {
		locationTimes = append(locationTimes, &LocationTime{
			Location: key,
			Time:     val,
		})
	}
	return locationTimes
}

func getLocationsFromQueries(queries url.Values) map[string]string {
	local, _ := time.LoadLocation("Local")

	locations := make(map[string]string)
	locations["Local"] = getCurTime(local)

	var err error
	var location *time.Location
	if queries.Has("location") {
		reqLocations := queries["location"]
		for _, reqLocation := range reqLocations {
			if location, err = time.LoadLocation(reqLocation); err != nil {
				locations[reqLocation] = err.Error()
			} else {
				locations[reqLocation] = getCurTime(location)
			}
		}
	}

	return locations
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
