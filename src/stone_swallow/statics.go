package stone_swallow

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/runtime"
	"google.golang.org/appengine/taskqueue"
)

// StatTotal is Datastoreのすべてのエンティティ
type StatTotal struct {
	Count                 int64
	Bytes                 int64
	Entity_bytes          int64
	Builtin_index_bytes   int64
	Builtin_index_count   int64
	Composite_index_bytes int64
	Composite_index_count int64
	Timestamp             time.Time
}

func (entity *StatTotal) Load(p []datastore.Property) error {
	err := datastore.LoadStruct(entity, p)
	if fmerr, ok := err.(*datastore.ErrFieldMismatch); ok && fmerr != nil && fmerr.Reason == "no such struct field" {
		// ignore
	} else if err != nil {
		return err
	}

	return nil
}

func (entity *StatTotal) Save() ([]datastore.Property, error) {
	p, err := datastore.SaveStruct(entity)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getStatics(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	results := map[string]interface{}{}

	// Task Queue Stats 取得
	qsl, err := taskqueue.QueueStats(ctx, []string{"default"})
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(qsl)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "TaskQueue.QueueStats:%s", string(b))
	results["TaskQueue.QueueStats"] = qsl

	// Runtime 取得
	stats, err := runtime.Stats(ctx)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err = json.Marshal(stats)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "Runtime.Stats:%s", string(b))
	results["Runtime.Stats"] = stats

	// Datastore Stats 取得
	var e []StatTotal
	_, err = datastore.NewQuery("__Stat_Total__").GetAll(ctx, &e)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	results["Datastore.Stats"] = e

	// HTTP ResponseのためのJsonに変換する
	b, err = json.Marshal(results)
	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
