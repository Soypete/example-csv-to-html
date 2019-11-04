package main

import (
	"encoding/csv"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"
)

type Speaker struct {
	Time      string
	Name      string
	TalkTitle string
	Abstract  string
}

func main() {
	f, err := os.Open("lightning_talks.csv")
	if err != nil {
		log.Fatal("file error")
	}
	reader := csv.NewReader(f)
	file, err := reader.ReadAll()
	if err != nil {
		log.Fatal("error in reader")
	}
	var speakers []Speaker
	for i, f := range file {
		if i == 0 {
			continue
		}
		speakers = append(speakers, Speaker{
			Name:      f[2],
			TalkTitle: f[1],
			Abstract:  f[3],
		})
	}
	// fmt.Println(speakers)
	speakers = randomize(speakers)
	speakers = setTime(speakers)
	// fmt.Println(speakers)
	t := template.New("t").Funcs(templateFuncs)
	t, err = t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = t.Execute(w, speakers)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func setTime(s []Speaker) []Speaker {
	for i := range s {
		s[i].Time = time.Date(2019, time.November, 5, 18, 30+i*13, 0, 0, time.UTC).Format(time.Kitchen)
	}
	return s
}

var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

// In the template, we use rangeStruct to turn our struct values
// into a slice we can iterate over
var htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<style>
table, th, td {
	border: 1px solid black;
  }
</style>  
<body>
<h2>Lightning Talk Schedule</h2>
<table style="width:100%">
<tr>
<th>Time</th>
<th>Name</th>
<th>Talk Title</th>
<th>Abstract</th>
</tr>
{{range .}}<tr>
{{range rangeStruct .}} <td>{{.}}</td>
{{end}}</tr>
{{end}}
</table>
</body>
</html>`

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func randomize(vals []Speaker) []Speaker {
	// r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, v := range vals {
		n := len(vals)
		randIndex := rand.Intn(n - 1)
		s := vals[randIndex]
		vals[randIndex] = v
		vals[i] = s
	}
	return vals
}
