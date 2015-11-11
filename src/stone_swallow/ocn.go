package stone_swallow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

type OCNMessage struct {
	Kind           string
	Id             string
	SelfLink       string
	Name           string
	Bucket         string
	Generation     string
	Metageneration string
	ContentType    string
	Updated        time.Time
	StrageClass    string
	Size           string
	Md5Hash        string
	MediaLink      string
	Owner          ACL
	Crc32c         string
	Etag           string
}

type ACL struct {
	Entity   string
	EntityId string
}

func handlerOCNReceiver(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	for k, v := range r.Header {
		log.Infof(ctx, "%s:%s", k, v)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "ERROR request body read: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, string(body))

	if r.Header.Get("X-Goog-Resource-State") == "sync" {
		log.Infof(ctx, "sync message")
		w.WriteHeader(http.StatusOK)
		return
	}

	var m OCNMessage
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&m)
	if err != nil {
		log.Errorf(ctx, "ERROR json decode: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("X-Goog-Resource-State") == "exists" {
		_, err = taskqueue.Add(ctx,
			&taskqueue.Task{
				Path: fmt.Sprintf("/tq/1/importBilling?bucket=%s&fileName=%s", m.Bucket, m.Name),
			},
			"billingimport")
		if err != nil {
			log.Errorf(ctx, "ERROR billingimport task add: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("done!"))
}
