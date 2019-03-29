package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

//Index holds fields displayed on the index.html template
type Index struct {
	AppName          string
	AppInstanceIndex int
	Database         string
	SpaceName        string
}

func main() {

	index := Index{"Unknown", -1, "Unknown", "Unknown"}

	template := template.Must(template.ParseFiles("templates/index.html"))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		appEnv, _ := cfenv.Current()
		if appEnv.Name != "" {
			index.AppName = appEnv.Name
		}
		if appEnv.Index > -1 {
			index.AppInstanceIndex = appEnv.Index
		}
		if appEnv.SpaceName != "" {
			index.SpaceName = appEnv.SpaceName
		}

		for k := range appEnv.Services {
			index.Database = k
		}

		if err := template.ExecuteTemplate(w, "index.html", index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(1)
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
