package stone_swallow

import (
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"

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

func createHogeKey(c context.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Hoge", id, 0, nil)
}

func (f *Hoge) key(c context.Context) *datastore.Key {
	return createHogeKey(c, f.Id)
}

func (f *Hoge) save(c context.Context) (*Hoge, error) {
	f.Created = time.Now()
	k, err := datastore.Put(c, f.key(c), f)
	if err != nil {
		return nil, err
	}
	f.Id = k.StringID()
	return f, nil
}

func postLog(w http.ResponseWriter, r *http.Request, c context.Context) {
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	log.Infof(c, body)
}

func putHoge(w http.ResponseWriter, r *http.Request, c context.Context) {
	id := r.FormValue("id")

	var h Hoge
	h.Id = id
	h.Name = "グフ"
	_, err := h.save(c)
	if err != nil {
		log.Errorf(c, "hoge save error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(h)
	if err != nil {
		log.Errorf(c, "hoge json encode error: %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHoge(c context.Context, key *datastore.Key) (Hoge, error) {
	var hoge Hoge
	err := datastore.Get(c, key, &hoge)
	if err != nil {
		return hoge, err
	}
	return hoge, nil
}

func getAllHoge(c context.Context) ([]Hoge, error) {
	hoges := []Hoge{}
	_, err := datastore.NewQuery("Hoge").GetAll(c, &hoges)
	if err != nil {
		return nil, err
	}
	return hoges, nil
}
