package stone_swallow

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func getStatics(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
