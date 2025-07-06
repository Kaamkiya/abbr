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

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		to := r.Form.Get("to")
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
		w.Write([]byte(key))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
