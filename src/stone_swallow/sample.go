package stone_swallow

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Hoge struct {
	Id      string    `json:"id" datastore:"-"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

func createHogeKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Hoge", id, 0, nil)
}

func (f *Hoge) key(c appengine.Context) *datastore.Key {
	return createHogeKey(c, f.Id)
}

func (f *Hoge) save(c appengine.Context) (*Hoge, error) {
	f.Created = time.Now()
	k, err := datastore.Put(c, f.key(c), f)
	if err != nil {
		return nil, err
	}
	f.Id = k.StringID()
	return f, nil
}

func postLog(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	c.Infof(body)
}

func putHoge(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	id := r.FormValue("id")

	var h Hoge
	h.Id = id
	h.Name = "グフ"
	_, err := h.save(c)
	if err != nil {
		c.Errorf("hoge save error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		c.Errorf("hoge json encode error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHoge(c appengine.Context, key *datastore.Key) (Hoge, error) {
	var hoge Hoge
	err := datastore.Get(c, key, &hoge)
	if err != nil {
		return hoge, err
	}
	return hoge, nil
}

func getAllHoge(c appengine.Context) ([]Hoge, error) {
	hoges := []Hoge{}
	_, err := datastore.NewQuery("Hoge").GetAll(c, &hoges)
	if err != nil {
		return nil, err
	}
	return hoges, nil
}
