package main

import (
	"fmt"
	"net/http"

	"github.com/Kaamkiya/nanoid-go"
)

var pairs = map[string]string{}

func main() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)

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
		html := r.Form.Get("webform") == "true"
		add := true
		key := nanoid.Default()

		for k, urlTo := range pairs {
			if urlTo == to {
				add = false
				key = k
			}
		}

		if add {
			pairs[key] = to
		}

		if html {
			w.Write([]byte(`<!DOCTYPE html><p>Short URL: <a href="/` + key + `">` + key + `</a></p>`))
			return
		}

		w.Write([]byte(key))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)

		path := r.URL.Path
		if path[0] == '/' {
			path = path[1:]
		}

		to := pairs[path]

		if to != "" {
			http.Redirect(w, r, to, http.StatusPermanentRedirect)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	http.ListenAndServe(":4000", nil)
}
