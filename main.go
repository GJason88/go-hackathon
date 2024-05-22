package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type LocationTime struct {
	Location string
	Time     string
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "Hello World!")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	queries := r.URL.Query()

	locations := getLocationsFromQueries(queries)
	locationTimes := buildLocationTimes(locations)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locationTimes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildLocationTimes(locations map[string]string) []*LocationTime {
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

func timeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	queries := r.URL.Query()

	locations := getLocationsFromQueries(queries)
	locationTimes := buildLocationTimes(locations)
	timeTemplate := showTimes(locationTimes, getAllLocations())

	fmt.Fprint(w, timeTemplate)
}

func getAllLocations() []string {
	var zoneDirs = []string{
		// Update path according to your OS
		"/usr/share/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/usr/lib/locale/TZ/",
	}
	var allLocations []string
	for _, zoneDir := range zoneDirs {
        ReadFile("", zoneDir, &allLocations)
    }
	return allLocations
}

func ReadFile(path string, zoneDir string, allLocations *[]string) {
    files, _ := os.ReadDir(zoneDir + path)
    for _, f := range files {
        if f.Name() != strings.ToUpper(f.Name()[:1]) + f.Name()[1:] {
            continue
        }
        if f.IsDir() {
            ReadFile(path + "/" + f.Name(), zoneDir, allLocations)
        } else {
            *allLocations = append(*allLocations, (path + "/" + f.Name())[1:])
        }
    }
}

func showTimes(timeList []*LocationTime, allLocations []string) string {
	type PageData struct {
		TimeList []*LocationTime
		AllLocations []string
	}
	timeData := PageData{ TimeList: timeList, AllLocations: allLocations }
	
	const html5 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>TimeZone-a-rama</title>
	<style>
	html { height:100%; background:#f8f8f8; font-size:21px }
@font-face {
    font-family:'Source Sans Pro';
    font-weight:400;
    font-style:normal;
    font-stretch:normal;
    src:url('/style/lib/fonts/source-sans-pro/SourceSansPro-Regular.ttf.woff2') format('woff2'),
         url('/style/lib/fonts/source-sans-pro/SourceSansPro-Regular.otf.woff') format('woff'),
         url('/style/lib/fonts/source-sans-pro/SourceSansPro-Regular.otf') format('opentype'),
         url('/style/lib/fonts/source-sans-pro/SourceSansPro-Regular.ttf') format('truetype');
}
body { padding:10%; margin:0; min-height:100%; position:relative;
	background:transparent; font-family:"Source Sans Pro",helvetica,sans-serif }
	h1,h3,h4 { margin: 0.2rem 0; padding:0 }
			</style>
  </head>
  <body>
	<h1>TimeZone and Times Browser</h1>
	<form style="float:left; width: 35%; border:1px solid gray; border-radius:0.5rem;padding:1rem 0 2rem 1em; box-shadow: 0.2rem 0.1rem 1rem 0.2rem rgba(0, 10, 30, 0.2);">
	<select name="location" size="15" multiple>
	{{range $_, $loc := .AllLocations }}
	<option value="{{ $loc }}">{{ $loc }}</option>
	{{else}}
	<option value="">NOPE</option>
	{{end}}
	</select>
	<div style="padding-left:2rem">
	<input type="submit" value="Update Time Display">
	</div>
	</form>
	<div style="float:left; margin-left:1em; width: 45%; border:1px solid gray; border-radius:0.5rem;padding:0.5rem 0 1rem 1em; box-shadow: 0.2rem 0.1rem 1rem 0.2rem rgba(0, 10, 30, 0.2);">
		{{range $_, $tz := .TimeList}}
		<div>
		<h3 style="font-weight:normal">{{ $tz.Location }}</h3>
		<h4 style="padding-left:1em; opacity:0.6;">{{ $tz.Time }}</h4>
		</div>
		{{else}}
		<div><strong>No times!</strong></div>
		{{end}}
	</div>
  </body>
</html>
`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	t, err := template.New("webpage").Parse(html5)
	check(err)

	b := new(bytes.Buffer)
	err = t.Execute(b, timeData)
	check(err)
	return b.String()
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/time/template", timeTemplateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
