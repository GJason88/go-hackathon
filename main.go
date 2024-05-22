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
	queries := r.URL.Query()

	locations := getLocationsFromQueries(queries)
	locationTimes := buildLocationTimes(locations)
	timeTemplate := showTimes(locationTimes, make([]string, 0))

	fmt.Fprint(w, timeTemplate)
}

func getAllLocations() []string {
	var zoneDirs = []string{
		// Update path according to your OS
		"/usr/share/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/usr/lib/locale/TZ/",
	}
	for _, zoneDir := range zoneDirs {
        ReadFile("", zoneDir)
    }
}

func ReadFile(path string, zoneDir string, allLocations []string) {
    files, _ := os.ReadDir(zoneDir + path)
    for _, f := range files {
        if f.Name() != strings.ToUpper(f.Name()[:1]) + f.Name()[1:] {
            continue
        }
        if f.IsDir() {
            ReadFile(path + "/" + f.Name(), zoneDir)
        } else {
            append(allLocations, (path + "/" + f.Name())[1:])
        }
    }
}

func showTimes(timeList []*LocationTime, allLocations []string) string {
	type PageData struct {
		timeList *LocationTime
		location string
	}
	
	timeData := PageData{ timeList, allLocations }
	
	const html5 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>TimeZone-a-rama</title>
	<style>
	html { height:100%; background:#f8f8f8; font-size:16px }
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
body { padding:0; margin:0; min-height:100%; position:relative;
	background:transparent; font-family:"Source Sans Pro",helvetica,sans-serif }
		</style>
  </head>
  <body>
	<h1>TimeZone Browser</h1>
	<form>
	<select name="location" size="15">
	{{range .allLocations }}
	<option value="{{ . }}">{{ . }}</option>
	{{else}}
	<option value="">NOPE</option>
	{{end}}
	</select>
	</form>
		{{range .timeLlist}}
		<div>{{ .Time }} {{ .Location }}</div>
		{{else}}
		<div><strong>No times!</strong></div>
		{{end}}
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
//	{ tl => timeList, tz => timeZone }
	err = t.Execute(b, ARGUMENTS)
	check(err)
	return b.String()
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/time/template", timeTemplateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
