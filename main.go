package main

import (
	"encoding/json"
	"errors"
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
		log.Printf("unexpected error updating active component while reading request body: %v", bserr)
		// Treat read errors as network errors.
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	var update activeComponentUpdate
	var syntaxError *json.SyntaxError
	if err := json.Unmarshal(bs, &update); errors.As(err, &syntaxError) {
		http.Error(w, "request body is not valid JSON", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("unexpected error unmarshaling request body while updating active component: %v", err)
		// Treat other unmarshal errors similarly to syntax errors.
		http.Error(w, "request body is not valid JSON", http.StatusBadRequest)
		return
	}
	activeComponent = update.Component
	relativePath = update.RelativePath
}

func root(w http.ResponseWriter, r *http.Request) {
	t, terr := template.ParseFiles(`C:\Users\Icosatess\Source\minimapui\index.html`)
	if terr != nil {
		log.Printf("Failed to parse templates in root handler: %v", terr)
		// The template was missing or invalid, probably.
		http.Error(w, "Missing files to render minimap. See logs.", http.StatusInternalServerError)
		return
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
		log.Printf("Failed to marshal active component update: %v", bserr)
		http.Error(w, "Couldn't generate active component update", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(bs); err != nil {
		log.Printf("Got error writing active component response, ignoring: %v", err)
		return
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
