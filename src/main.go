package stone_swallow

import (
	"net/http"
	"encoding/json"
	"appengine"
	"appengine/datastore"
	"time"
	"log"
)

type Hoge struct {
	Id          string    `json:"id" datastore:"-"`
	Created     time.Time `json:"created"`
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/post", postData)
}

func (f *Hoge) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Hoge", f.Id, 0, nil)
}

func (f *Hoge) save(c appengine.Context) (*Hoge, error) {
	f.Created = time.Now()
	k, err := datastore.Put(c, f.Key(c), f)
	if err != nil {
		return nil, err
	}
	f.Id = k.StringID()
	return f, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	kind := r.FormValue("kind")

	var dst []datastore.PropertyList
	q := datastore.NewQuery(kind)
	k, err := q.GetAll(c, &dst);
	if err != nil {
		c.Errorf("handler error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dst)
}

func postData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	log.Printf("id=%s", "Ohhhh")

	var h Hoge
	h.Id = "gufu"
	h.save(c)
	err := json.NewEncoder(w).Encode(h)
	if err != nil {
		c.Errorf("postData error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
