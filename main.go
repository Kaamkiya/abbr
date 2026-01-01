package main

import (
	"net/http"

	"github.com/Kaamkiya/nanoid-go"
)

var pairs = map[string]string{}

func main() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Pattern != "/create" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != "POST" && r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Method == "GET" {
			http.ServeFile(w, r, "create.html")
			return
		}

		r.ParseForm()
		to := r.Form.Get("to")
		isHTML := r.Form.Get("webform") == "true"
		shouldAdd := true
		key := nanoid.Default()

		for k, urlTo := range pairs {
			if urlTo == to {
				shouldAdd = false
				key = k
			}
		}

		if shouldAdd {
			pairs[key] = to
		}

		if isHTML {
			w.Write([]byte(`<!DOCTYPE html><p>Short URL: http://localhost:4000/` + key + `</p>`))
			return
		}

		w.Write([]byte(key))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" {
			http.Redirect(w, r, "/create", http.StatusPermanentRedirect)
			return
		}

		path = path[1:]

		to := pairs[path]

		if to != "" {
			http.Redirect(w, r, to, http.StatusPermanentRedirect)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	http.ListenAndServe(":4000", nil)
}
