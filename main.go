package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

// activeComponent describes which component (VS Code workspace folder) is
// currently active. Empty string means no component is active.
var activeComponent string

func updateActiveComponent(w http.ResponseWriter, r *http.Request) {
	bs, bserr := io.ReadAll(r.Body)
	if bserr != nil {
		panic(bserr)
	}
	activeComponent = string(bs)
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

func main() {
	http.HandleFunc("POST /component/", updateActiveComponent)
	http.Handle("/public/", http.FileServer(http.Dir(`C:\Users\Icosatess\Source\minimapui`)))
	http.HandleFunc("/", root)

	srvAddr := "127.0.0.1:8081"
	log.Printf("Starting minimap server at %s", srvAddr)
	log.Fatal(http.ListenAndServe(srvAddr, nil))
}
