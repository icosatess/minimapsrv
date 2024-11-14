package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

// activeComponent describes which component (VS Code workspace folder) is
// currently active. Empty string means no component is active.
var activeComponent string
var relativePath string

type activeComponentUpdate struct {
	Component    string `json:"component"`
	RelativePath string `json:"relativePath"`
}

func updateActiveComponent(w http.ResponseWriter, r *http.Request) {
	bs, bserr := io.ReadAll(r.Body)
	if bserr != nil {
		panic(bserr)
	}
	var update activeComponentUpdate
	if err := json.Unmarshal(bs, &update); err != nil {
		panic(err)
	}
	activeComponent = update.Component
	relativePath = update.RelativePath
}

func root(w http.ResponseWriter, r *http.Request) {
	t, terr := template.ParseFiles(`C:\Users\Icosatess\Source\minimapui\index.html`)
	if terr != nil {
		panic(terr)
	}

	var componentString string
	if activeComponent != "" {
		componentString = "minimap-" + activeComponent
	}

	t.Execute(w, struct{ Component string }{
		Component: componentString,
	})
}

func getActiveComponent(w http.ResponseWriter, r *http.Request) {
	bs, bserr := json.Marshal(activeComponentUpdate{
		Component:    activeComponent,
		RelativePath: relativePath,
	})
	if bserr != nil {
		panic(bserr)
	}
	if _, err := w.Write(bs); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("POST /component/", updateActiveComponent)
	http.HandleFunc("GET /component/", getActiveComponent)
	http.Handle("/public/", http.FileServer(http.Dir(`C:\Users\Icosatess\Source\minimapui`)))
	http.HandleFunc("/", root)

	srvAddr := "127.0.0.1:8081"
	log.Printf("Starting minimap server at %s", srvAddr)
	log.Fatal(http.ListenAndServe(srvAddr, nil))
}
