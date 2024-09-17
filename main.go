package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

var activeComponent string = "server"

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
	t.Execute(w, struct{ Component string }{
		Component: "minimap-" + activeComponent,
	})
}

func main() {
	http.HandleFunc("POST /component/", updateActiveComponent)
	http.Handle("/public/", http.FileServer(http.Dir(`C:\Users\Icosatess\Source\minimapui`)))
	http.HandleFunc("/", root)

	log.Printf("Starting static file server at 127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
