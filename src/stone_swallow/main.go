package stone_swallow

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type runtimeEnv struct {
	NumCPU       int
	GOMAXPROCS   int
	NumGoroutine int
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	switch r.URL.Path {
	default:
		http.Error(w, "not found.", http.StatusNotFound)
	case "/entity":
		listEntity(w, r, c)
	case "/sample":
		putHoge(w, r, c)
	case "/env":
		listEnvironment(w, r, c)
	}
}

func listEnvironment(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))
	fmt.Println(runtime.NumGoroutine())

	re := &runtimeEnv{
		NumCPU:       runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS((0)),
		NumGoroutine: runtime.NumGoroutine(),
	}

	json, err := json.Marshal(re)
	if err != nil {
		c.Errorf("handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		c.Errorf("write response error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func listEntity(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	kind := r.FormValue("kind")
	log.Printf("kind=%s", kind)

	var dst []datastore.PropertyList
	q := datastore.NewQuery(kind)
	_, err := q.GetAll(c, &dst)
	if err != nil {
		c.Errorf("handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(dst)
	if err != nil {
		c.Errorf("handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(json)
	if err != nil {
		c.Errorf("write response error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
