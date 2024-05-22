package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	m := make(map[string]string)
	showTimes(m)
}

func showTimes(timeList map[string]string) {
	const html5 = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>TimeZone-a-rama</title>
  </head>
  <body>
		{{range .}}<div>{{ . }}</div>{{else}}<div><strong>No times!</strong></div>{{end}}
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

	err = t.Execute(os.Stdout, timeList)
	check(err)

}
